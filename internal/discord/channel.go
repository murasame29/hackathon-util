package discord

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

type Channel struct {
	ss *discordgo.Session
}

func newChannel(ss *discordgo.Session) *Channel {
	return &Channel{ss: ss}
}

var (
	ChannelTypeCategory = discordgo.ChannelTypeGuildCategory
	ChannelTypeText     = discordgo.ChannelTypeGuildText
	ChannelTypeVoice    = discordgo.ChannelTypeGuildVoice
)

func (c *Channel) Create(ctx context.Context, category discordgo.ChannelType, guildID, name, parentID string) (*discordgo.Channel, error) {
	var (
		result *discordgo.Channel
		err    error
	)
	if category == ChannelTypeCategory {
		result, err = c.ss.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
			Name: name,
			Type: category,
		})
	} else {
		result, err = c.ss.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
			Name:     name,
			Type:     category,
			ParentID: parentID,
		})
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}
