package discordgo

import (
	"github.com/bwmarrin/discordgo"
)

type DiscordSession struct {
	ss *discordgo.Session
}

func New(ss *discordgo.Session) *DiscordSession {
	return &DiscordSession{ss: ss}
}
