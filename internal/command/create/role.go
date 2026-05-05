package create

import (
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// mentorRoleColor is the Discord role color for mentors (#3498db, Belize Hole blue).
const mentorRoleColor = 3447003

// createParticipantsRole ensures the participants and mentor roles exist.
func createParticipantsRole(dg *discordgo.Session, guildID, eventName string, existingRoles map[string]string, mentionable bool, dryRun bool) (string, string, error) {
	participantsRoleID := ensureRole(dg, guildID, "参加者_"+eventName, existingRoles, &discordgo.RoleParams{
		Mentionable: &mentionable,
	}, dryRun)
	color := mentorRoleColor
	mentorRoleID := ensureRole(dg, guildID, "メンター_"+eventName, existingRoles, &discordgo.RoleParams{
		Mentionable: &mentionable,
		Color:       &color,
	}, dryRun)
	return participantsRoleID, mentorRoleID, nil
}

// ensureRole returns the existing role ID or creates a new one.
// In dry-run mode the API call is skipped and an empty string is returned for new roles.
func ensureRole(dg *discordgo.Session, guildID, roleName string, existingRoles map[string]string, params *discordgo.RoleParams, dryRun bool) string {
	if id := existingRoles[roleName]; id != "" {
		slog.Info("role already exists, skipping", slog.String("role", roleName))
		return id
	}
	if dryRun {
		slog.Info("dry run: would create role", slog.String("role", roleName))
		return ""
	}
	params.Name = roleName
	role, err := dg.GuildRoleCreate(guildID, params)
	if err != nil {
		slog.Error("failed to create role", slog.String("role", roleName), slog.String("error.message", err.Error()))
		return ""
	}
	slog.Info("role created", slog.String("role", roleName))
	return role.ID
}

// assignRoleIfMissing assigns roleID to the user if they don't already have it.
// In dry-run mode the member is still fetched but the role assignment is skipped.
func assignRoleIfMissing(dg *discordgo.Session, guildID, userID, roleID, username string, dryRun bool) error {
	member, err := dg.GuildMember(guildID, userID)
	if err != nil || member == nil {
		slog.Info("could not retrieve member, skipping", slog.String("username", username), slog.String("role_id", roleID))
		return nil
	}
	for _, r := range member.Roles {
		if r == roleID {
			slog.Info("member already has role, skipping", slog.String("username", username), slog.String("role_id", roleID))
			return nil
		}
	}
	if dryRun {
		slog.Info("dry run: would assign role", slog.String("role_id", roleID), slog.String("username", username))
		return nil
	}
	if err := dg.GuildMemberRoleAdd(guildID, userID, roleID); err != nil {
		return err
	}
	slog.Info("role assigned", slog.String("role_id", roleID), slog.String("username", username))
	return nil
}

// assignRoleToRowMembers assigns roleID to every member listed in columns B–F.
// Returns usernames that could not be found in userMap.
func assignRoleToRowMembers(gs *guildState, row []any, roleID string, userMap map[string]string, dryRun bool) []string {
	var notFound []string
	for i := 1; i <= 5; i++ {
		username := columnUsername(row, i)
		if username == "" {
			continue
		}
		userID, ok := userMap[username]
		if !ok {
			slog.Info("username not found in guild, skipping", slog.String("username", username), slog.String("role_id", roleID))
			notFound = append(notFound, username)
			continue
		}
		if err := assignRoleIfMissing(gs.dg, gs.guildID, userID, roleID, username, dryRun); err != nil {
			slog.Error("failed to assign role", slog.String("role_id", roleID), slog.String("username", username), slog.String("error.message", err.Error()))
		}
	}
	return notFound
}

func columnUsername(row []any, i int) string {
	if i >= len(row) {
		return ""
	}
	s, ok := row[i].(string)
	if !ok {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(s))
}

// rowUsernames extracts non-empty usernames from columns B–F (index 1–5).
func rowUsernames(row []any) []string {
	names := make([]string, 0, 5)
	for i := 1; i <= 5; i++ {
		if u := columnUsername(row, i); u != "" {
			names = append(names, u)
		}
	}
	return names
}
