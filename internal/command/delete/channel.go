package delete

import (
	"log/slog"
	"time"

	"github.com/bwmarrin/discordgo"
)

func deleteTeamCategoryAndChildren(dg *discordgo.Session, guildID, categoryName string, dryRun bool) {
	channels, err := dg.GuildChannels(guildID)
	if err != nil {
		slog.Error("fetching channels failed", slog.String("error.message", err.Error()))
		return
	}

	categoryID := findCategoryID(channels, categoryName)
	if categoryID == "" {
		slog.Info("category not found, skipping", slog.String("category", categoryName))
		return
	}

	deleteChildChannels(dg, channels, categoryID, dryRun)
	deleteSingleChannel(dg, categoryID, categoryName, "category", dryRun)
}

// findCategoryID returns the ID of the named category, or "".
func findCategoryID(channels []*discordgo.Channel, categoryName string) string {
	for _, ch := range channels {
		if ch.Type == discordgo.ChannelTypeGuildCategory && ch.Name == categoryName {
			return ch.ID
		}
	}
	return ""
}

// deleteChildChannels deletes all channels whose parent is categoryID.
func deleteChildChannels(dg *discordgo.Session, channels []*discordgo.Channel, categoryID string, dryRun bool) {
	for _, ch := range channels {
		if ch.ParentID != categoryID {
			continue
		}
		deleteSingleChannel(dg, ch.ID, ch.Name, "channel", dryRun)
		if !dryRun {
			time.Sleep(250 * time.Millisecond)
		}
	}
}

// deleteSingleChannel deletes one channel (or logs a dry-run message).
func deleteSingleChannel(dg *discordgo.Session, id, name, kind string, dryRun bool) {
	if dryRun {
		slog.Info("dry run: would delete "+kind, slog.String(kind, name), slog.String(kind+"_id", id))
		return
	}
	if _, err := dg.ChannelDelete(id); err != nil {
		slog.Error("deleting "+kind+" failed", slog.String(kind, name), slog.String("error.message", err.Error()))
	} else {
		slog.Info(kind+" deleted", slog.String(kind, name))
	}
}
