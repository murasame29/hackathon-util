package discordgo

import (
	"github.com/bwmarrin/discordgo"
	"github.com/murasame29/hackathon-util/cmd/config"
)

type DiscordSession struct {
	ss *discordgo.Session
}

func New() *DiscordSession {
	dg, err := discordgo.New("Bot " + config.Config.Discord.BotToken)
	if err != nil {
		panic(err)
	}

	return &DiscordSession{ss: dg}
}
