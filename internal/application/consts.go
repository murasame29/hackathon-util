package application

import (
	"errors"
)

var (
	ErrNoSetDiscordBotToken = errors.New("no set discord bot token")
	ErrNoSetDiscordGuildID  = errors.New("no set discord GuildID")
	ErrNoSetDataSource      = errors.New("no set data source")
)

type DataSourceMode string

const (
	DataSourceModeGoogleSheet DataSourceMode = "google_sheets"
	DataSourceModeFile        DataSourceMode = "file"
)
