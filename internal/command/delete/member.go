package delete

import (
	"log/slog"
	"time"

	"github.com/bwmarrin/discordgo"
)

func removeRoleFromAllMembers(dg *discordgo.Session, guildID, roleID, roleName string, dryRun bool) {
	slog.Info("removing role from all holders", slog.String("role", roleName))
	after := ""
	removed := 0
	for {
		members, err := dg.GuildMembers(guildID, after, 1000)
		if err != nil {
			slog.Error("fetching members failed", slog.String("error.message", err.Error()))
			break
		}
		if len(members) == 0 {
			break
		}
		removed += processMemberRoleRemoval(dg, guildID, members, roleID, roleName, dryRun, &after)
		time.Sleep(200 * time.Millisecond)
	}
	slog.Info("role removals completed", slog.String("role", roleName), slog.Int("removed", removed))
}

// processMemberRoleRemoval iterates one page of members and removes the role.
// Updates after in-place and returns the count of removals performed.
func processMemberRoleRemoval(
	dg *discordgo.Session,
	guildID string,
	members []*discordgo.Member,
	roleID, roleName string,
	dryRun bool,
	after *string,
) int {
	removed := 0
	for _, m := range members {
		*after = m.User.ID
		if !memberHasRole(m, roleID) {
			continue
		}
		if dryRun {
			slog.Info("dry run: would remove role", slog.String("role", roleName), slog.String("username", m.User.Username))
			continue
		}
		if err := dg.GuildMemberRoleRemove(guildID, m.User.ID, roleID); err != nil {
			slog.Error("removing role failed", slog.String("role", roleName), slog.String("username", m.User.Username), slog.String("error.message", err.Error()))
		} else {
			removed++
		}
	}
	return removed
}

func memberHasRole(m *discordgo.Member, roleID string) bool {
	for _, r := range m.Roles {
		if r == roleID {
			return true
		}
	}
	return false
}
