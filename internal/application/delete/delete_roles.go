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

type DeleteRolesOptions struct {
	SheetID     string
	Range       string
	FilePath    string
	EnvFilePath string

	config *config.EnvironmentsVariables

	dataSourceMode application.DataSourceMode
}

func NewDeleteRolesOptions() *DeleteRolesOptions {
	return &DeleteRolesOptions{
		config: config.NewEnvironment(),
	}
}

func (o *DeleteRolesOptions) Complete() error {
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

func (o *DeleteRolesOptions) Validate() error {
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

func (o *DeleteRolesOptions) Run() error {
	discord, err := discord.NewDiscord(o.config.Discord.BotToken)
	if err != nil {
		return err
	}

	cr := newDeleteRole(discord)
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

	return cr.Execute(result.Teams)
}

type deleteRole struct {
	discord *discord.Discord
}

func newDeleteRole(discord *discord.Discord) *deleteRole {
	return &deleteRole{
		discord: discord,
	}
}

func (c *deleteRole) Execute(teams map[string][]string) error {
	var eg errgroup.Group
	ctx := context.Background()

	for teamName, _ := range teams {
		eg.Go(func() error {
			role, err := c.discord.Role.Get(ctx, config.Config.Discord.GuildID, teamName)
			if err != nil {
				return err
			}

			if err := c.discord.Role.Delete(ctx, config.Config.Discord.GuildID, role.ID); err != nil {
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
