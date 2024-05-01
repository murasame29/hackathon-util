package discordgo

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/murasame29/hackathon-util/cmd/config"
)

type DiscordSession struct {
	ss *discordgo.Session
}

func New() (*DiscordSession, error) {
	dg, err := discordgo.New("Bot " + config.Config.Discord.BotToken)
	if err != nil {
		return nil, fmt.Errorf("error creating Discord session %v", err)
	}

	return &DiscordSession{ss: dg}, nil
}
