package application

import (
	"github.com/murasame29/hackathon-util/internal/adapter/gateways/discordgo"
	"github.com/murasame29/hackathon-util/internal/adapter/gateways/gs"
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

func NewSheetLessApplicationService(ds *discordgo.DiscordSession) *ApplicationService {
	return &ApplicationService{
		ds: ds,
	}
}
