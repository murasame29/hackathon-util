package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/murasame29/hackathon-util/internal/application"
)

var integerOptionMinValue = 1.0
var BreakoutRoomCommand = &discordgo.ApplicationCommand{
	Name:        "breakout-room",
	Description: "start breakout room",
	Options: []*discordgo.ApplicationCommandOption{

		{
			Type:        discordgo.ApplicationCommandOptionChannel,
			Name:        "target-vc",
			Description: "Host VC channel name",
			ChannelTypes: []discordgo.ChannelType{
				discordgo.ChannelTypeGuildVoice,
				discordgo.ChannelTypeGuildStageVoice,
			},
			Required: true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "number-per-room",
			Description: "Number of people in each room",
			MinValue:    &integerOptionMinValue,
			MaxValue:    100,
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "time-to-live",
			Description: "Breakout room time-to-live",
			MaxValue:    120,
			Required:    true,
		},
	},
}

func (dh *DiscordHandler) BreakoutRoom(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	// This example stores the provided arguments in an []interface{}
	// which will be used to format the bot's response
	margs := make([]interface{}, 0, len(options))

	// Get the value from the option map.
	// When the option exists, ok = true
	if option, ok := optionMap["target-vc"]; ok {
		// Option values must be type asserted from interface{}.
		// Discordgo provides utility functions to make this simple.
		margs = append(margs, option.ChannelValue)
	}

	if option, ok := optionMap["number-per-room"]; ok {
		margs = append(margs, option.IntValue())
	}

	if option, ok := optionMap["time-to-live"]; ok {
		margs = append(margs, option.IntValue())
	}

	breakoutRoomParams := application.BreakoutRoomParams{
		TargetVC:      margs[0].(discordgo.Channel),
		NumberPerRoom: margs[1].(int),
		TimeToLive:    margs[2].(int),
	}
	// TODO
	app := application.NewSheetLessApplicationService()
	app.BreakoutRoomStart(breakoutRoomParams)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "health ok!",
		},
	})
}
