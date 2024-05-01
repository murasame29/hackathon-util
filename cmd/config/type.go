package config

var Config *config

type Env string

const (
	Dev  Env = "dev"
	Prod Env = "prod"
)

type config struct {
	Application struct {
		Env Env `env:"ENV"`
	}

	Spreadsheets struct {
		ID    string `env:"GOOGLE_SPREADSHEET_ID"`
		Range string `env:"GOOGLE_SPREADSHEET_RANGE"`
	}

	Discord struct {
		BotToken string `env:"DISCORD_BOT_TOKEN"`
		GuildID  string `env:"DISCORD_GUILD_ID"`
	}
}
