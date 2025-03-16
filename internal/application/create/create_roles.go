package create

import (
	"context"
	"os"

	"github.com/murasame29/hackathon-util/cmd/config"
	"github.com/murasame29/hackathon-util/internal/application"
	"github.com/murasame29/hackathon-util/internal/discord"
	"golang.org/x/sync/errgroup"
)

type CreateRolesOptions struct {
	SheetID  string
	Range    string
	FilePath string

	config *config.EnvironmentsVariables

	dataSourceMode application.DataSourceMode
}

func NewCreateRolesOptions() *CreateRolesOptions {
	return &CreateRolesOptions{
		config: config.NewEnvironment(),
	}
}

func (o *CreateRolesOptions) Complete() error {
	o.config.Discord.BotToken = os.Getenv("DISCORD_BOT_TOKEN")
	o.config.Discord.GuildID = os.Getenv("DISCORD_GUILD_ID")
	return nil
}

func (o *CreateRolesOptions) Validate() error {
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

func (o *CreateRolesOptions) Run() error {
	discord, err := discord.NewDiscord(o.config.Discord.BotToken)
	if err != nil {
		return err
	}

	cr := newCreateRole(discord)
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

	return cr.Execute(result.Teams)
}

type createRole struct {
	discord *discord.Discord
}

func newCreateRole(discord *discord.Discord) *createRole {
	return &createRole{
		discord: discord,
	}
}

func (c *createRole) Execute(teams map[string][]string) error {
	var eg errgroup.Group
	ctx := context.Background()

	for teamName, _ := range teams {
		eg.Go(func() error {
			_, err := c.discord.Role.Create(ctx, config.Config.Discord.GuildID, teamName)
			if err != nil {
				return err
			}
			return nil
		})

	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}
