package apply

import (
	"os"

	"github.com/murasame29/hackathon-util/cmd/config"
	"github.com/murasame29/hackathon-util/internal/application"
	"github.com/murasame29/hackathon-util/internal/application/bind"
	"github.com/murasame29/hackathon-util/internal/application/create"
	"github.com/murasame29/hackathon-util/internal/datasource"
	"github.com/murasame29/hackathon-util/internal/datasource/csv"
	"github.com/murasame29/hackathon-util/internal/datasource/sheet"
	"github.com/murasame29/hackathon-util/internal/discord"
)

type ApplyOptions struct {
	SheetID     string
	Range       string
	FilePath    string
	EnvFilePath string

	config *config.EnvironmentsVariables

	dataSourceMode application.DataSourceMode
}

func NewApplyOptions() *ApplyOptions {
	return &ApplyOptions{
		config: config.NewEnvironment(),
	}
}

func (o *ApplyOptions) Complete() error {
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

func (o *ApplyOptions) Validate() error {
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

func (o *ApplyOptions) Run() error {
	discord, err := discord.NewDiscord(o.config.Discord.BotToken)
	if err != nil {
		return err
	}

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

	if err := create.NewCreateChannel(discord).Execute(result.TeamNames); err != nil {
		return err
	}

	if err := create.NewCreateRole(discord).Execute(result.Teams); err != nil {
		return err
	}

	return bind.NewBindRole(discord).Execute(result.Teams)
}
