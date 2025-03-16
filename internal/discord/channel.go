package discord

import (
	"context"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/murasame29/hackathon-util/pkg/cache"
)

type Channel struct {
	ss *discordgo.Session

	cache *cache.Cache
}

func newChannel(ss *discordgo.Session) *Channel {
	return &Channel{ss: ss, cache: cache.New()}
}

var (
	ChannelTypeCategory = discordgo.ChannelTypeGuildCategory
	ChannelTypeText     = discordgo.ChannelTypeGuildText
	ChannelTypeVoice    = discordgo.ChannelTypeGuildVoice
)

func (c *Channel) Get(ctx context.Context, category discordgo.ChannelType, guildID, name string) (*discordgo.Channel, error) {
	if channel, ok := c.cache.Get(fmt.Sprintf("%s_%d_%s", guildID, category, name)); ok {
		return channel.(*discordgo.Channel), nil
	}

	channels, err := c.ss.GuildChannels(guildID, discordgo.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	for _, channel := range channels {
		c.cache.Set(fmt.Sprintf("%s_%d_%s", guildID, channel.Type, channel.Name), channel)
		if channel.Name == name {
			return channel, nil
		}
	}

	return nil, ErrResourceNotFound
}

func (c *Channel) Create(ctx context.Context, category discordgo.ChannelType, guildID, name, parentID string) (*discordgo.Channel, error) {
	var (
		result *discordgo.Channel
		err    error
	)

	ok, err := c.Exist(ctx, category, guildID, name)
	if err != nil {
		return nil, err
	}

	if ok {
		return nil, ErrResourceAlreadyExists
	}

	if category == ChannelTypeCategory {
		result, err = c.ss.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
			Name: name,
			Type: category,
		}, discordgo.WithContext(ctx))
	} else {
		result, err = c.ss.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
			Name:     name,
			Type:     category,
			ParentID: parentID,
		}, discordgo.WithContext(ctx))
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Channel) Exist(ctx context.Context, category discordgo.ChannelType, guildID, name string) (bool, error) {
	_, err := c.Get(ctx, category, guildID, name)
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
