package application

import "errors"

var (
	ErrNoSetDiscordBotToken = errors.New("no set discord bot token")
	ErrNoSetDiscordGuildID  = errors.New("no set discord GuildID")
	ErrNoSetDataSource      = errors.New("no set data source")
)

type DataSourceMode string

const (
	DataSourceURL  DataSourceMode = "url"
	DataSourceFile DataSourceMode = "file"
)
