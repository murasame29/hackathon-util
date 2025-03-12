package discord

import "github.com/bwmarrin/discordgo"

type User struct {
	ss *discordgo.Session
}

func newUser(ss *discordgo.Session) *User {
	return &User{ss: ss}
}
