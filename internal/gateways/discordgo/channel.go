package discordgo

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/cockroachdb/errors"
)

type CreateChannelParams struct {
	GuildID string
}

func (ds *DiscordSession) CreateChannelCategory(ctx context.Context, serverID string, name string) (string, error) {
	channel, err := ds.ss.GuildChannelCreate(serverID, name, discordgo.ChannelTypeGuildCategory, discordgo.WithContext(ctx))
	if err != nil {
		return "", errors.Wrap(err, "Failed to create channel category")
	}

	return channel.ID, nil
}

func (ds *DiscordSession) CreateChannelText(ctx context.Context, serverID, categoryID, name string) (string, error) {
	channel, err := ds.ss.GuildChannelCreateComplex(serverID, discordgo.GuildChannelCreateData{
		Name:     name,
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: categoryID,
	}, discordgo.WithContext(ctx))
	if err != nil {
		return "", errors.Wrap(err, "Failed to create channel text")
	}

	return channel.ID, nil
}

func (ds *DiscordSession) CreateChannelVoice(ctx context.Context, serverID, categoryID, name string) (string, error) {
	channel, err := ds.ss.GuildChannelCreateComplex(serverID, discordgo.GuildChannelCreateData{
		Name:      name,
		ParentID:  categoryID,
		Type:      discordgo.ChannelTypeGuildVoice,
		Bitrate:   64000,
		UserLimit: 20,
	}, discordgo.WithContext(ctx))
	if err != nil {
		return "", errors.Wrap(err, "Failed to create channel voice")
	}

	return channel.ID, nil
}

func (ds *DiscordSession) MoveCategory(ctx context.Context, categoryID, channelID string) error {
	_, err := ds.ss.ChannelEdit(channelID, &discordgo.ChannelEdit{ParentID: categoryID}, discordgo.WithContext(ctx))
	if err != nil {
		return errors.Wrap(err, "Failed to move channel")
	}

	return nil
}

func (ds *DiscordSession) GetChannel(ctx context.Context, serverID string) ([]*discordgo.Channel, error) {
	channels, err := ds.ss.GuildChannels(serverID, discordgo.WithContext(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get channels")
	}

	return channels, nil
}

func (ds *DiscordSession) DeleteChannel(ctx context.Context, channelID string) error {
	_, err := ds.ss.ChannelDelete(channelID, discordgo.WithContext(ctx))
	if err != nil {
		return err
	}
	return nil
}
