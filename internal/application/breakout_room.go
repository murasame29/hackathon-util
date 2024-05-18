package application

import (
	"github.com/bwmarrin/discordgo"
)

type BreakoutRoomParams struct {
	TargetVC      discordgo.Channel
	NumberPerRoom int
	TimeToLive    int
}

func (as *ApplicationService) BreakoutRoomStart(params BreakoutRoomParams) (string, error) {

	return "message", nil
}
