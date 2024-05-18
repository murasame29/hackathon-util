package discord

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/murasame29/hackathon-util/cmd/config"
	"github.com/murasame29/hackathon-util/pkg/logger"
)

type DiscordHandler struct {
	ss *discordgo.Session
}

func NewHandler(ss *discordgo.Session) *DiscordHandler {
	return &DiscordHandler{
		ss: ss,
	}
}

func (dh *DiscordHandler) Open(ctx context.Context) error {
	logger.Info(ctx, "discord bot opend")
	if err := dh.ss.Open(); err != nil {
		logger.Error(ctx, err.Error())
		return err
	}

	// add handler
	dh.implCommand()
	dh.createCommand()
	//
	return nil
}

func (dh *DiscordHandler) createCommand() {
	commands := []*discordgo.ApplicationCommand{
		{
			Name: "health",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "health check commnad",
		},
	}
	id := dh.ss.State.User.ID
	guildIDs := config.Config.Discord.GuildID
	for _, command := range commands {
		for _, guildID := range guildIDs {
			_, err := dh.ss.ApplicationCommandCreate(id, guildID, command)
			logger.Info(context.Background(), command.Name)
			if err != nil {
				logger.Error(context.Background(), err.Error())
			}
		}
	}
}

func (dh *DiscordHandler) implCommand() {
	implCommand := map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		"health": dh.Health,
	}

	dh.ss.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if cmd, ok := implCommand[i.ApplicationCommandData().Name]; ok {
			cmd(s, i)
		}
	})
}
