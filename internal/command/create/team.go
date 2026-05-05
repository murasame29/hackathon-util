package create

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/murasame29/hackathon-util/internal/config"
)

// processTeamRow handles one row from the spreadsheet.
// Returns usernames that could not be found in the guild.
func processTeamRow(
	gs *guildState,
	cfg *config.Config,
	row []any,
	teamName, mentorRoleID, participantsRoleID string,
	baseVCOverwrites []*discordgo.PermissionOverwrite,
	mentionable bool,
	dryRun bool,
) []string {
	roleID, ok := ensureTeamRole(gs, teamName, mentionable, dryRun)
	if !ok {
		return nil
	}

	vcOverwrites, categoryOverwrites := buildTeamOverwrites(cfg, gs.guildID, roleID, mentorRoleID, baseVCOverwrites)

	if !ensureTeamCategory(gs, cfg, teamName, categoryOverwrites, vcOverwrites, dryRun) {
		return nil
	}

	usernames := rowUsernames(row)
	userMap, err := lookupUserIDs(gs.dg, gs.guildID, usernames)
	if err != nil {
		slog.Error("failed to lookup guild members", slog.String("error.message", err.Error()))
		return nil
	}

	teamMissing := assignRoleToRowMembers(gs, row, roleID, userMap, dryRun)
	participantsMissing := assignRoleToRowMembers(gs, row, participantsRoleID, userMap, dryRun)
	return append(teamMissing, participantsMissing...)
}

// ensureTeamRole returns the role ID for teamName, creating it if necessary.
// Returns ("", false) on error.
func ensureTeamRole(gs *guildState, teamName string, mentionable bool, dryRun bool) (string, bool) {
	if id, exists := gs.existingRoles[teamName]; exists {
		slog.Info("role already exists, skipping", slog.String("team", teamName))
		return id, true
	}
	if dryRun {
		slog.Info("dry run: would create role", slog.String("team", teamName))
		return "", true
	}
	role, err := gs.dg.GuildRoleCreate(gs.guildID, &discordgo.RoleParams{
		Name:        teamName,
		Mentionable: &mentionable,
	})
	if err != nil {
		slog.Error("role create failed", slog.String("team", teamName), slog.String("error.message", err.Error()))
		return "", false
	}
	gs.existingRoles[teamName] = role.ID
	slog.Info("role created", slog.String("team", teamName))
	return role.ID, true
}

// buildTeamOverwrites returns (vcOverwrites, categoryOverwrites) based on privacy settings.
func buildTeamOverwrites(
	cfg *config.Config,
	guildID, roleID, mentorRoleID string,
	baseVCOverwrites []*discordgo.PermissionOverwrite,
) ([]*discordgo.PermissionOverwrite, []*discordgo.PermissionOverwrite) {
	categoryOverwrites := buildPublicPermissionOverwrites(guildID)
	vcOverwrites := baseVCOverwrites
	if cfg.Discord.EnablePrivateCategory {
		categoryOverwrites = buildCategoryPermissionOverwrites(roleID, mentorRoleID, guildID, cfg.Discord.MuteRoleID)
		vcOverwrites = categoryOverwrites
	}
	return vcOverwrites, categoryOverwrites
}
