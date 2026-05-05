package create

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/murasame29/hackathon-util/internal/config"
	"github.com/murasame29/hackathon-util/internal/googlesheet"
)

// Config holds runtime options for the create command.
type Config struct {
	BotToken string
	Cfg      *config.Config
}

// guildState holds pre-fetched Discord guild data shared across team processing.
type guildState struct {
	dg                 *discordgo.Session
	guildID            string
	existingRoles      map[string]string
	existingCategories map[string]string
	channels           []*discordgo.Channel
}

// Run executes the sheet-to-discord channel/role creation flow.
func Run(cc Config) error {
	cfg := cc.Cfg

	teamData, err := googlesheet.GetTeamData(cfg.GoogleSheet.ID, cfg.GoogleSheet.TeamTableRange, cfg.GoogleSheet.CredentialFile)
	if err != nil {
		return fmt.Errorf("failed to fetch team data: %w", err)
	}

	dg, err := discordgo.New("Bot " + cc.BotToken)
	if err != nil {
		return fmt.Errorf("failed to create Discord session: %w", err)
	}
	defer dg.Close()

	gs, err := loadGuildState(dg, cfg.Discord.GuildID)
	if err != nil {
		return err
	}

	baseVCOverwrites, participantsRoleID, mentorRoleID := buildBaseOverwrites(dg, cfg, gs)

	notFoundUsers := processAllTeams(gs, cfg, teamData, mentorRoleID, participantsRoleID, baseVCOverwrites)
	logNotFoundUsers(notFoundUsers)
	return nil
}

// buildBaseOverwrites creates the participants/mentor roles and returns the base VC permission overwrites.
func buildBaseOverwrites(dg *discordgo.Session, cfg *config.Config, gs *guildState) ([]*discordgo.PermissionOverwrite, string, string) {
	mentionable := true
	participantsRoleID, mentorRoleID, _ := createParticipantsRole(dg, cfg.Discord.GuildID, cfg.EventName, gs.existingRoles, mentionable)

	baseVCOverwrites := buildPublicPermissionOverwrites(cfg.Discord.GuildID)
	if cfg.Discord.EnablePrivateVC {
		baseVCOverwrites = buildVCPermissionOverwrites(participantsRoleID, mentorRoleID, cfg.Discord.GuildID)
	}
	return baseVCOverwrites, participantsRoleID, mentorRoleID
}

// processAllTeams iterates over sheet rows and processes each team.
// Returns all usernames that could not be found in the guild.
func processAllTeams(
	gs *guildState,
	cfg *config.Config,
	teamData [][]any,
	mentorRoleID, participantsRoleID string,
	baseVCOverwrites []*discordgo.PermissionOverwrite,
) []string {
	mentionable := true
	var notFoundUsers []string
	for _, row := range teamData {
		if len(row) == 0 {
			continue
		}
		teamName := fmt.Sprintf("%v", row[0])
		if teamName == "" {
			continue
		}
		missing := processTeamRow(gs, cfg, row, teamName, mentorRoleID, participantsRoleID, baseVCOverwrites, mentionable)
		notFoundUsers = append(notFoundUsers, missing...)
		slog.Info("team processed", slog.String("team", teamName))
	}
	return notFoundUsers
}

// logNotFoundUsers logs usernames that could not be resolved in the guild.
func logNotFoundUsers(notFoundUsers []string) {
	for _, name := range notFoundUsers {
		slog.Warn("user not found in guild", slog.String("username", name))
	}
}

// loadGuildState fetches roles, channels, and categories from Discord.
func loadGuildState(dg *discordgo.Session, guildID string) (*guildState, error) {
	roles, err := dg.GuildRoles(guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch roles: %w", err)
	}
	existingRoles := make(map[string]string, len(roles))
	for _, r := range roles {
		existingRoles[r.Name] = r.ID
	}

	channels, err := dg.GuildChannels(guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch channels: %w", err)
	}
	existingCategories := make(map[string]string)
	for _, ch := range channels {
		if ch.Type == discordgo.ChannelTypeGuildCategory {
			existingCategories[ch.Name] = ch.ID
		}
	}

	return &guildState{
		dg:                 dg,
		guildID:            guildID,
		existingRoles:      existingRoles,
		existingCategories: existingCategories,
		channels:           channels,
	}, nil
}
