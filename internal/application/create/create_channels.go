package create

import (
	"context"
	"fmt"
	"os"

	"github.com/murasame29/hackathon-util/cmd/config"
	"github.com/murasame29/hackathon-util/internal/application"
	"github.com/murasame29/hackathon-util/internal/discord"
	"golang.org/x/sync/errgroup"
)

type CreateChannelsOptions struct {
	SheetID  string
	Range    string
	FilePath string

	config *config.EnvironmentsVariables

	dataSourceMode application.DataSourceMode
}

func NewCreateChannelsOptions() *CreateChannelsOptions {
	return &CreateChannelsOptions{}
}

func (o *CreateChannelsOptions) Complete() error {
	o.config.Discord.BotToken = os.Getenv("DISCORD_BOT_TOKEN")
	o.config.Discord.GuildID = os.Getenv("DISCORD_GUILD_ID")
	return nil
}

func (o *CreateChannelsOptions) Validate() error {
	if o.config.Discord.BotToken == "" {
		return application.ErrNoSetDiscordBotToken
	}

	if o.config.Discord.GuildID == "" {
		return application.ErrNoSetDiscordGuildID
	}

	if (o.SheetID == "" || o.Range == "") && o.FilePath == "" {
		return application.ErrNoSetDataSource
	}

	o.dataSourceMode = application.DataSourceModeFile
	if o.FilePath == "" {
		o.dataSourceMode = application.DataSourceModeGoogleSheet
	}

	return nil
}

func (o *CreateChannelsOptions) Run() error {
	discord, err := discord.NewDiscord(o.config.Discord.BotToken)
	if err != nil {
		return err
	}

	cc := newCreateChannel(discord)
	var result *application.ReadDataSourceResult

	switch o.dataSourceMode {
	case application.DataSourceModeFile:
		result, err = application.NewDataSourceCSV(o.FilePath).Read()
	case application.DataSourceModeGoogleSheet:
		result, err = application.NewDataSourceGoogleSheets(o.SheetID, o.Range).Read()
	}

	if err != nil {
		return err
	}

	return cc.Execute(result.TeamNames)
}

type createChannel struct {
	discord *discord.Discord
}

func newCreateChannel(discord *discord.Discord) *createChannel {
	return &createChannel{
		discord: discord,
	}
}

func (c *createChannel) Execute(categories []string) error {
	var eg errgroup.Group
	ctx := context.Background()

	for _, category := range categories {
		eg.Go(func() error {
			category, err := c.discord.Channel.Create(ctx, discord.ChannelTypeCategory, config.Config.Discord.GuildID, category, "")
			if err != nil {
				return fmt.Errorf("failed to create category: %w", err)
			}

			_, err = c.discord.Channel.Create(ctx, discord.ChannelTypeText, config.Config.Discord.GuildID, "雑談", category.ID)
			if err != nil {
				return fmt.Errorf("failed to create text channel 雑談: %w", err)
			}

			_, err = c.discord.Channel.Create(ctx, discord.ChannelTypeText, config.Config.Discord.GuildID, "会議", category.ID)
			if err != nil {
				return fmt.Errorf("failed to create text channel 会議: %w", err)
			}

			_, err = c.discord.Channel.Create(ctx, discord.ChannelTypeVoice, config.Config.Discord.GuildID, "vc", category.ID)
			if err != nil {
				return fmt.Errorf("failed to create voice channel vc: %w", err)
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}
