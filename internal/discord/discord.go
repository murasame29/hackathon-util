package discord

import "github.com/bwmarrin/discordgo"

type Discord struct {
	Channel *Channel
	Role    *Role
	User    *User
}

func NewDiscord(token string) (*Discord, error) {
	ss, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &Discord{
		Channel: newChannel(ss),
		Role:    newRole(ss),
		User:    newUser(ss),
	}, nil
}
