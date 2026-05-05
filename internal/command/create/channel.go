package create

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/murasame29/hackathon-util/internal/config"
)

// ensureTeamCategory upserts the category and its text/voice children.
// Returns false if the category could not be created.
func ensureTeamCategory(
	gs *guildState,
	cfg *config.Config,
	teamName string,
	categoryOverwrites, vcOverwrites []*discordgo.PermissionOverwrite,
	dryRun bool,
) bool {
	categoryID, ok := upsertCategory(gs, cfg.Discord.GuildID, teamName, categoryOverwrites, dryRun)
	if !ok {
		return false
	}
	upsertChildChannel(gs, cfg.Discord.GuildID, teamName, categoryID, "やりとり", discordgo.ChannelTypeGuildText, categoryOverwrites, dryRun)
	upsertChildChannel(gs, cfg.Discord.GuildID, teamName, categoryID, "会話", discordgo.ChannelTypeGuildVoice, vcOverwrites, dryRun)
	return true
}

// upsertCategory creates the category if absent, otherwise updates its permissions.
// Returns (categoryID, true) on success. In dry-run mode writes are skipped.
func upsertCategory(gs *guildState, guildID, teamName string, overwrites []*discordgo.PermissionOverwrite, dryRun bool) (string, bool) {
	if id, exists := gs.existingCategories[teamName]; exists {
		if dryRun {
			slog.Info("dry run: would update category permission", slog.String("team", teamName))
		} else if err := updateChannelPermissions(gs.dg, id, overwrites); err != nil {
			slog.Error("update category permission failed", slog.String("team", teamName), slog.String("error.message", err.Error()))
		} else {
			slog.Info("category permission updated", slog.String("team", teamName))
		}
		return id, true
	}
	if dryRun {
		slog.Info("dry run: would create category", slog.String("team", teamName))
		return "", true
	}
	ch, err := gs.dg.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
		Name:                 teamName,
		Type:                 discordgo.ChannelTypeGuildCategory,
		PermissionOverwrites: overwrites,
	})
	if err != nil {
		slog.Error("category create failed", slog.String("team", teamName), slog.String("error.message", err.Error()))
		return "", false
	}
	gs.existingCategories[teamName] = ch.ID
	slog.Info("category created", slog.String("team", teamName))
	return ch.ID, true
}

// upsertChildChannel updates permissions if the channel exists, otherwise creates it.
// In dry-run mode writes are skipped.
func upsertChildChannel(
	gs *guildState,
	guildID, teamName, categoryID, name string,
	chType discordgo.ChannelType,
	overwrites []*discordgo.PermissionOverwrite,
	dryRun bool,
) {
	if id := findChildChannelID(gs.channels, categoryID, name, chType); id != "" {
		if dryRun {
			slog.Info("dry run: would update channel permission", slog.String("team", teamName), slog.String("channel", name))
		} else if err := updateChannelPermissions(gs.dg, id, overwrites); err != nil {
			slog.Error("update channel permission failed", slog.String("team", teamName), slog.String("channel", name), slog.String("error.message", err.Error()))
		} else {
			slog.Info("channel permission updated", slog.String("team", teamName), slog.String("channel", name))
		}
		return
	}
	if dryRun {
		slog.Info("dry run: would create channel", slog.String("team", teamName), slog.String("channel", name))
		return
	}
	if _, err := gs.dg.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
		Name:                 name,
		Type:                 chType,
		ParentID:             categoryID,
		PermissionOverwrites: overwrites,
	}); err != nil {
		slog.Error("channel create failed", slog.String("team", teamName), slog.String("channel", name), slog.String("error.message", err.Error()))
	} else {
		slog.Info("channel created", slog.String("team", teamName), slog.String("channel", name))
	}
}

// findChildChannelID returns the ID of a child channel matching name and type, or "".
func findChildChannelID(channels []*discordgo.Channel, categoryID, name string, chType discordgo.ChannelType) string {
	for _, ch := range channels {
		if ch.ParentID == categoryID && ch.Name == name && ch.Type == chType {
			return ch.ID
		}
	}
	return ""
}

func updateChannelPermissions(dg *discordgo.Session, channelID string, overwrites []*discordgo.PermissionOverwrite) error {
	_, err := dg.ChannelEditComplex(channelID, &discordgo.ChannelEdit{
		PermissionOverwrites: overwrites,
	})
	return err
}
