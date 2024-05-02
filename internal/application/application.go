package application

import (
	"github.com/murasame29/hackathon-util/internal/gateways/discordgo"
	"github.com/murasame29/hackathon-util/internal/gateways/gs"
)

type ApplicationService struct {
	gs *gs.GoogleSpreadSeet
	ds *discordgo.DiscordSession
}

func NewApplicationService(gs *gs.GoogleSpreadSeet, ds *discordgo.DiscordSession) *ApplicationService {
	return &ApplicationService{
		gs: gs,
		ds: ds,
	}
}
