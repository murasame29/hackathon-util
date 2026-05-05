package delete

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/murasame29/hackathon-util/internal/config"
	"github.com/murasame29/hackathon-util/internal/googlesheet"
)

// Config holds runtime options for the delete command.
type Config struct {
	BotToken         string
	DryRun           bool
	RemoveAllMembers bool
	Cfg              *config.Config
}

// Run executes the sheet-to-discord-delete cleanup flow.
func Run(dc Config) error {
	cfg := dc.Cfg
	slog.Info("starting delete", slog.Bool("dry_run", dc.DryRun), slog.Bool("remove_all_members", dc.RemoveAllMembers))

	teamData, err := googlesheet.GetTeamData(cfg.GoogleSheet.ID, cfg.GoogleSheet.TeamTableRange, cfg.GoogleSheet.CredentialFile)
	if err != nil {
		return fmt.Errorf("failed to fetch team data: %w", err)
	}

	dg, err := discordgo.New("Bot " + dc.BotToken)
	if err != nil {
		return fmt.Errorf("failed to create Discord session: %w", err)
	}
	defer dg.Close()

	existingRoles, err := fetchExistingRoles(dg, cfg.Discord.GuildID)
	if err != nil {
		return err
	}

	teamNames := collectTeamNames(teamData)
	if len(teamNames) == 0 {
		slog.Warn("no team names found in the sheet; nothing to delete")
		return nil
	}

	deleteTeams(dg, cfg.Discord.GuildID, teamNames, existingRoles, dc.DryRun)
	removeParticipantsRoleIfNeeded(dg, cfg, existingRoles, dc)
	printDeleteSummary(dc.DryRun)
	return nil
}

// fetchExistingRoles returns a name→ID map of all roles in the guild.
func fetchExistingRoles(dg *discordgo.Session, guildID string) (map[string]string, error) {
	roles, err := dg.GuildRoles(guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch roles: %w", err)
	}
	m := make(map[string]string, len(roles))
	for _, r := range roles {
		m[r.Name] = r.ID
	}
	return m, nil
}

// collectTeamNames extracts unique team names from the first column of sheet data.
func collectTeamNames(teamData [][]any) map[string]struct{} {
	names := make(map[string]struct{})
	for _, row := range teamData {
		if len(row) == 0 {
			continue
		}
		name := strings.TrimSpace(fmt.Sprintf("%v", row[0]))
		if name != "" {
			names[name] = struct{}{}
		}
	}
	return names
}

// deleteTeams deletes categories and roles for every team.
func deleteTeams(dg *discordgo.Session, guildID string, teamNames map[string]struct{}, existingRoles map[string]string, dryRun bool) {
	for name := range teamNames {
		deleteTeamCategoryAndChildren(dg, guildID, name, dryRun)
	}
	for name := range teamNames {
		deleteTeamRole(dg, guildID, name, existingRoles, dryRun)
		time.Sleep(200 * time.Millisecond)
	}
}

// deleteTeamRole removes all members from a role then deletes it.
func deleteTeamRole(dg *discordgo.Session, guildID, name string, existingRoles map[string]string, dryRun bool) {
	roleID, ok := existingRoles[name]
	if !ok {
		slog.Info("role not found for team, skipping", slog.String("team", name))
		return
	}
	removeRoleFromAllMembers(dg, guildID, roleID, name, dryRun)
	if dryRun {
		slog.Info("dry run: would delete role", slog.String("role", name), slog.String("role_id", roleID))
		return
	}
	if err := dg.GuildRoleDelete(guildID, roleID); err != nil {
		slog.Error("deleting role failed", slog.String("role", name), slog.String("error.message", err.Error()))
	} else {
		slog.Info("role deleted", slog.String("role", name))
	}
}

// removeParticipantsRoleIfNeeded strips the participants role from all members when requested.
func removeParticipantsRoleIfNeeded(dg *discordgo.Session, cfg *config.Config, existingRoles map[string]string, dc Config) {
	if !dc.RemoveAllMembers {
		return
	}
	roleName := "参加者_" + cfg.EventName
	roleID, ok := existingRoles[roleName]
	if !ok {
		slog.Info("participants role not found, skipping", slog.String("role", roleName))
		return
	}
	removeRoleFromAllMembers(dg, cfg.Discord.GuildID, roleID, roleName, dc.DryRun)
}

func printDeleteSummary(dryRun bool) {
	if dryRun {
		fmt.Println("\n🧪 DRY RUN complete. Set --dry-run=false to perform actual deletions.")
	} else {
		fmt.Println("\n🧹 Cleanup complete.")
	}
}
