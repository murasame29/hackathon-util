package driver

import (
	"github.com/bwmarrin/discordgo"
	"github.com/murasame29/hackathon-util/cmd/config"
)

func NewDiscordSession() *discordgo.Session {
	dg, err := discordgo.New("Bot " + config.Config.Discord.BotToken)
	if err != nil {
		panic(err)
	}

	return dg
}
