package discord

import "github.com/bwmarrin/discordgo"

type Role struct {
	ss *discordgo.Session
}

func newRole(ss *discordgo.Session) *Role {
	return &Role{ss: ss}
}
