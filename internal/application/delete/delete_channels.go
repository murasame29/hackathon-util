package delete

import (
	"context"
	"os"

	"github.com/murasame29/hackathon-util/cmd/config"
	"github.com/murasame29/hackathon-util/internal/application"
	"github.com/murasame29/hackathon-util/internal/datasource"
	"github.com/murasame29/hackathon-util/internal/datasource/csv"
	"github.com/murasame29/hackathon-util/internal/datasource/sheet"
	"github.com/murasame29/hackathon-util/internal/discord"
	"golang.org/x/sync/errgroup"
)

type DeleteChannelsOptions struct {
	SheetID     string
	Range       string
	FilePath    string
	EnvFilePath string

	config *config.EnvironmentsVariables

	dataSourceMode application.DataSourceMode
}

func NewDeleteChannelsOptions() *DeleteChannelsOptions {
	return &DeleteChannelsOptions{
		config: config.NewEnvironment(),
	}
}

func (o *DeleteChannelsOptions) Complete() error {
	o.config.Discord.BotToken = os.Getenv("DISCORD_BOT_TOKEN")
	o.config.Discord.GuildID = os.Getenv("DISCORD_GUILD_ID")

	if o.EnvFilePath != "" {
		if err := config.LoadEnv(o.EnvFilePath); err != nil {
			return err
		}

		o.config = config.Config
	}
	return nil
}

func (o *DeleteChannelsOptions) Validate() error {
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

func (o *DeleteChannelsOptions) Run() error {
	discord, err := discord.NewDiscord(o.config.Discord.BotToken)
	if err != nil {
		return err
	}

	cc := newDeleteChannel(discord)
	var result *datasource.ReadDataSourceResult

	switch o.dataSourceMode {
	case application.DataSourceModeFile:
		result, err = csv.NewDataSource(o.FilePath).Read()
	case application.DataSourceModeGoogleSheet:
		result, err = sheet.NewDataSource(o.SheetID, o.Range).Read()
	}

	if err != nil {
		return err
	}

	return cc.Execute(result.TeamNames)
}

type deleteChannel struct {
	discord *discord.Discord
}

func newDeleteChannel(discord *discord.Discord) *deleteChannel {
	return &deleteChannel{
		discord: discord,
	}
}

func (c *deleteChannel) Execute(categories []string) error {
	var eg errgroup.Group
	ctx := context.Background()

	for _, category := range categories {
		eg.Go(func() error {
			category, err := c.discord.Channel.Get(ctx, discord.ChannelTypeCategory, config.Config.Discord.GuildID, category)
			if err != nil {
				return err
			}

			zatudanChannel, err := c.discord.Channel.Get(ctx, discord.ChannelTypeText, config.Config.Discord.GuildID, "雑談", discord.WithParentID(category.ID))
			if err != nil {
				return err
			}

			kaigiChannel, err := c.discord.Channel.Get(ctx, discord.ChannelTypeText, config.Config.Discord.GuildID, "会議", discord.WithParentID(category.ID))
			if err != nil {
				return err
			}

			voiceChannel, err := c.discord.Channel.Get(ctx, discord.ChannelTypeVoice, config.Config.Discord.GuildID, "vc", discord.WithParentID(category.ID))
			if err != nil {
				return err
			}

			deleteChannels := []string{
				zatudanChannel.ID,
				kaigiChannel.ID,
				voiceChannel.ID,
				category.ID,
			}

			for _, channelID := range deleteChannels {
				if err := c.discord.Channel.Delete(ctx, channelID); err != nil {
					return err
				}
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}
