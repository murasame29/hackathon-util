package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Discord struct {
	Channel *Channel
	Role    *Role
	Member  *Member
}

var (
	ErrResourceNotFound      = fmt.Errorf("resource not found")
	ErrResourceAlreadyExists = fmt.Errorf("resource already exists")
)

func NewDiscord(token string) (*Discord, error) {
	ss, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &Discord{
		Channel: newChannel(ss),
		Role:    newRole(ss),
		Member:  newMember(ss),
	}, nil
}
