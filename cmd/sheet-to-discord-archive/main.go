package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const defaultForumID = "1366003202331512842"

func mustEnv(key string) string {
	v := os.Getenv(key)
	if strings.TrimSpace(v) == "" {
		log.Fatalf("missing required env: %s", key)
	}
	return v
}

func getenv(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

func getenvBool(key string, def bool) bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	if v == "true" || v == "1" || v == "yes" {
		return true
	}
	if v == "false" || v == "0" || v == "no" {
		return false
	}
	return def
}

func splitCSV(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func getTeamNames(spreadsheetID, rangeStr, credentialsFile string) ([]string, error) {
	srv, err := sheets.NewService(context.Background(), option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, err
	}
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, rangeStr).Do()
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	var names []string
	for _, row := range resp.Values {
		if len(row) == 0 {
			continue
		}
		n := strings.TrimSpace(fmt.Sprintf("%v", row[0]))
		if n != "" && !seen[n] {
			seen[n] = true
			names = append(names, n)
		}
	}
	return names, nil
}

// sourceChannels returns (channel, categoryName) pairs to archive.
type sourceChannel struct {
	ch       *discordgo.Channel
	category string
}

func resolveSources(channels []*discordgo.Channel, teamNames, channelNames, explicitIDs []string) []sourceChannel {
	byID := map[string]*discordgo.Channel{}
	for _, c := range channels {
		byID[c.ID] = c
	}
	var out []sourceChannel
	if len(explicitIDs) > 0 {
		for _, id := range explicitIDs {
			c, ok := byID[id]
			if !ok {
				log.Printf("[SKIP] Channel not found: %s", id)
				continue
			}
			cat := ""
			if p, ok := byID[c.ParentID]; ok {
				cat = p.Name
			}
			out = append(out, sourceChannel{c, cat})
		}
		return out
	}
	wanted := map[string]bool{}
	for _, n := range channelNames {
		wanted[n] = true
	}
	for _, team := range teamNames {
		catID := ""
		for _, c := range channels {
			if c.Type == discordgo.ChannelTypeGuildCategory && c.Name == team {
				catID = c.ID
				break
			}
		}
		if catID == "" {
			log.Printf("[SKIP] Category not found: %s", team)
			continue
		}
		for _, c := range channels {
			if c.ParentID == catID && c.Type == discordgo.ChannelTypeGuildText && wanted[c.Name] {
				out = append(out, sourceChannel{c, team})
			}
		}
	}
	return out
}

func main() {
	_ = godotenv.Load()
	botToken := mustEnv("DISCORD_BOT_TOKEN")
	guildID := mustEnv("DISCORD_GUILD_ID")
	forumID := getenv("ARCHIVE_FORUM_ID", defaultForumID)
	maxUpload, err := strconv.ParseInt(getenv("MAX_UPLOAD_BYTES", "10485760"), 10, 64)
	if err != nil || maxUpload <= 0 {
		log.Fatalf("invalid MAX_UPLOAD_BYTES")
	}
	cfg := config{
		DryRun:    getenvBool("DRY_RUN", true),
		Force:     getenvBool("FORCE", false),
		MaxUpload: maxUpload,
		StatePath: getenv("ARCHIVE_STATE_FILE", "archive-state.json"),
	}
	log.Printf("[INFO] DRY_RUN=%v FORCE=%v MAX_UPLOAD_BYTES=%d STATE=%s", cfg.DryRun, cfg.Force, cfg.MaxUpload, cfg.StatePath)

	var wh *webhookClient
	if !cfg.DryRun {
		id, token, err := parseWebhookURL(mustEnv("ARCHIVE_WEBHOOK_URL"))
		if err != nil {
			log.Fatalf("ARCHIVE_WEBHOOK_URL: %v", err)
		}
		wh = &webhookClient{id: id, token: token, http: &http.Client{Timeout: 120 * time.Second}, delay: 400 * time.Millisecond}
	}

	explicit := splitCSV(os.Getenv("CHANNEL_IDS"))
	var teams []string
	if len(explicit) == 0 {
		teams, err = getTeamNames(mustEnv("GOOGLE_SPREADSHEET_ID"), mustEnv("TEAM_RANGE"), mustEnv("GOOGLE_CREDENTIALS_FILE"))
		if err != nil {
			log.Fatalf("failed to fetch team data: %v", err)
		}
	}

	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("failed to create Discord session: %v", err)
	}
	channels, err := dg.GuildChannels(guildID)
	if err != nil {
		log.Fatalf("failed to fetch channels: %v", err)
	}
	active, err := dg.GuildThreadsActive(guildID)
	if err != nil {
		log.Fatalf("failed to fetch active threads: %v", err)
	}
	state, err := loadState(cfg.StatePath)
	if err != nil {
		log.Fatalf("%v", err)
	}
	rd := &reader{dg: dg, guildID: guildID, nameCache: map[string]string{}}
	a := &archiver{cfg: cfg, rd: rd, wh: wh, httpc: &http.Client{Timeout: 120 * time.Second}, state: state,
		posts: rd.forumPostIndex(forumID, active.Threads)}
	log.Printf("[INFO] forum has %d existing posts (active + archived public)", len(a.posts))

	sources := resolveSources(channels, teams, splitCSV(getenv("ARCHIVE_CHANNEL_NAMES", "やりとり")), explicit)
	for _, s := range sources {
		base := buildTitle(s.category, s.ch.Name)
		a.run(target{SourceID: s.ch.ID, Title: base})
		for _, th := range rd.childThreads(s.ch.ID, active.Threads) {
			a.run(target{SourceID: th.ID, Title: buildTitle(s.category, s.ch.Name, th.Name), ParentSourceID: s.ch.ID})
		}
	}

	s := a.sum
	if cfg.DryRun {
		fmt.Printf("\nDRY RUN complete: %d source channels. No webhook calls were made. Set DRY_RUN=false to archive.\n", len(sources))
		return
	}
	fmt.Printf("\nDone: created %d, resumed %d, skipped(exists) %d, skipped(done) %d, failed %d, verify mismatches %d, messages posted %d\n",
		s.Created, s.Resumed, s.SkippedExists, s.SkippedDone, s.Failed, s.Mismatch, s.MessagesPosted)
}
