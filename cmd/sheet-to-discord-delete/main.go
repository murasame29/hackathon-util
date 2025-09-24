package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func loadEnv() {
	_ = godotenv.Load() // optional; allow running without .env if env vars are set elsewhere
}

func getTeamData(spreadsheetID, rangeStr, credentialsFile string) ([][]interface{}, error) {
	srv, err := sheets.NewService(context.Background(), option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, err
	}

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, rangeStr).Do()
	if err != nil {
		return nil, err
	}

	return resp.Values, nil
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if strings.TrimSpace(v) == "" {
		log.Fatalf("missing required env: %s", key)
	}
	return v
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

func buildUsernameToIDMap(dg *discordgo.Session, guildID string, max int) (map[string]string, error) {
	userMap := make(map[string]string)
	after := ""
	total := 0
	for total < max {
		members, err := dg.GuildMembers(guildID, after, 1000)
		if err != nil {
			return nil, err
		}
		if len(members) == 0 {
			break
		}
		for _, m := range members {
			if m.User != nil && m.User.Username != "" {
				userMap[strings.ToLower(m.User.Username)] = m.User.ID
				after = m.User.ID
				total++
			}
		}
	}
	return userMap, nil
}

func removeRoleFromAllMembers(dg *discordgo.Session, guildID, roleID, roleName string, dryRun bool) {
	log.Printf("[STEP] Remove role '%s' from all holders", roleName)
	after := ""
	removed := 0
	for {
		members, err := dg.GuildMembers(guildID, after, 1000)
		if err != nil {
			log.Printf("[ERROR] fetching members: %v", err)
			break
		}
		if len(members) == 0 {
			break
		}
		for _, m := range members {
			after = m.User.ID
			for _, r := range m.Roles {
				if r == roleID {
					if dryRun {
						log.Printf("[DRY] would remove role '%s' from %s", roleName, m.User.Username)
					} else {
						if err := dg.GuildMemberRoleRemove(guildID, m.User.ID, roleID); err != nil {
							log.Printf("[ERROR] removing role '%s' from %s: %v", roleName, m.User.Username, err)
						} else {
							removed++
						}
					}
					break
				}
			}
		}
		// be gentle with rate limits
		time.Sleep(200 * time.Millisecond)
	}
	log.Printf("[OK] Role removals completed: %d members updated for '%s'", removed, roleName)
}

func deleteTeamCategoryAndChildren(dg *discordgo.Session, guildID, categoryName string, dryRun bool) {
	channels, err := dg.GuildChannels(guildID)
	if err != nil {
		log.Printf("[ERROR] fetching channels: %v", err)
		return
	}

	var categoryID string
	for _, ch := range channels {
		if ch.Type == discordgo.ChannelTypeGuildCategory && ch.Name == categoryName {
			categoryID = ch.ID
			break
		}
	}
	if categoryID == "" {
		log.Printf("[SKIP] Category not found: %s", categoryName)
		return
	}

	// Delete child channels first
	for _, ch := range channels {
		if ch.ParentID == categoryID { // child of this category
			if dryRun {
				log.Printf("[DRY] would delete channel: #%s (%s)", ch.Name, ch.ID)
			} else {
				if _, err := dg.ChannelDelete(ch.ID); err != nil {
					log.Printf("[ERROR] deleting channel #%s: %v", ch.Name, err)
				} else {
					log.Printf("[OK] Deleted channel: #%s", ch.Name)
				}
				// avoid hammering
				time.Sleep(250 * time.Millisecond)
			}
		}
	}

	// Delete the category itself
	if dryRun {
		log.Printf("[DRY] would delete category: %s (%s)", categoryName, categoryID)
	} else {
		if _, err := dg.ChannelDelete(categoryID); err != nil {
			log.Printf("[ERROR] deleting category %s: %v", categoryName, err)
		} else {
			log.Printf("[OK] Deleted category: %s", categoryName)
		}
	}
}

func main() {
	loadEnv()

	spreadsheetID := mustEnv("GOOGLE_SPREADSHEET_ID")
	botToken := mustEnv("DISCORD_BOT_TOKEN")
	guildID := mustEnv("DISCORD_GUILD_ID")
	credentialsFile := mustEnv("GOOGLE_CREDENTIALS_FILE")
	teamRange := mustEnv("TEAM_RANGE")
	participantsRoleName := "参加者_" + mustEnv("EVENT_NAME")
	dryRun := getenvBool("DRY_RUN", true)
	removeAllMembers := getenvBool("REMOVE_ALL_MEMBERS", false)

	log.Printf("[INFO] DRY_RUN=%v REMOVE_ALL_MEMBERS=%v", dryRun, removeAllMembers)

	teamData, err := getTeamData(spreadsheetID, teamRange, credentialsFile)
	if err != nil {
		log.Fatalf("failed to fetch team data: %v", err)
	}

	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("failed to create Discord session: %v", err)
	}
	defer dg.Close()

	roles, err := dg.GuildRoles(guildID)
	if err != nil {
		log.Fatalf("failed to fetch roles: %v", err)
	}
	existingRoles := make(map[string]string)
	for _, r := range roles {
		existingRoles[r.Name] = r.ID
	}

	// Build a unique set of team names from the sheet (first column)
	teamNames := make(map[string]struct{})
	for _, row := range teamData {
		if len(row) == 0 {
			continue
		}
		name := strings.TrimSpace(fmt.Sprintf("%v", row[0]))
		if name == "" {
			continue
		}
		teamNames[name] = struct{}{}
	}

	if len(teamNames) == 0 {
		log.Println("[WARN] No team names found in the sheet; nothing to delete")
		return
	}

	// 1) Delete channels/categories for each team
	for name := range teamNames {
		deleteTeamCategoryAndChildren(dg, guildID, name, dryRun)
	}

	// 2) Remove team roles from members and delete the roles
	for name := range teamNames {
		roleID, ok := existingRoles[name]
		if !ok {
			log.Printf("[SKIP] Role not found for team '%s'", name)
			continue
		}
		removeRoleFromAllMembers(dg, guildID, roleID, name, dryRun)
		if dryRun {
			log.Printf("[DRY] would delete role: %s (%s)", name, roleID)
		} else {
			if err := dg.GuildRoleDelete(guildID, roleID); err != nil {
				log.Printf("[ERROR] deleting role '%s': %v", name, err)
			} else {
				log.Printf("[OK] Deleted role: %s", name)
			}
		}
		// polite pacing
		time.Sleep(200 * time.Millisecond)
	}

	// 3) Optionally remove ALL_MEMBERS role from all assigned users (but DO NOT delete the role)
	if removeAllMembers {
		if roleID, ok := existingRoles[participantsRoleName]; ok {
			removeRoleFromAllMembers(dg, guildID, roleID, participantsRoleName, dryRun)
		} else {
			log.Printf("[SKIP] ALL_MEMBERS role '%s' not found", participantsRoleName)
		}
	}

	if dryRun {
		fmt.Println("\n🧪 DRY RUN complete. Set DRY_RUN=false to perform actual deletions.")
	} else {
		fmt.Println("\n🧹 Cleanup complete.")
	}
}
