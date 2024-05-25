package main

import (
	"context"
	"flag"
	// "github.com/joho/godotenv"
	"github.com/murasame29/hackathon-util/cmd/config"
	"github.com/murasame29/hackathon-util/internal/container"
	"github.com/murasame29/hackathon-util/internal/framewrok/discord"
	"github.com/murasame29/hackathon-util/internal/server"
	"github.com/murasame29/hackathon-util/pkg/logger"
	"log"
	"os"
	"strings"
)

type envFlag []string

func (e *envFlag) String() string {
	return strings.Join(*e, ",")
}

func (e *envFlag) Set(v string) error {
	*e = append(*e, v)
	return nil
}

func init() {
	// Usage: eg. go run main.go -e .env -e hoge.env -e fuga.env ...
	var envFile envFlag
	flag.Var(&envFile, "e", "path to .env file \n eg. -e .env -e another.env . ")
	flag.Parse()
	log.Println(envFile)
	if err := config.LoadEnv(envFile...); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	ctx := logger.NewLoggerWithContext(context.Background())

	handler := container.NewSheetLessContainer()
	var discordHandler *discord.DiscordHandler
	if err := container.Provide(func(dh *discord.DiscordHandler) {
		discordHandler = dh
	}); err != nil {
		logger.Error(ctx, err.Error())
		return err
	}
	server.
		New(config.Config.Application.Addr, handler, discordHandler).
		RunWithGraceful(ctx)

	return nil
}
