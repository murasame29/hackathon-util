package main

import (
	"context"
	"flag"

	// "github.com/joho/godotenv"
	"log"
	"os"
	"strings"

	"github.com/murasame29/hackathon-util/cmd/config"
	"github.com/murasame29/hackathon-util/internal/container"
	"github.com/murasame29/hackathon-util/internal/framewrok/discord"
	"github.com/murasame29/hackathon-util/pkg/logger"
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

	if err := container.NewContainer(); err != nil {
		return err
	}

	discordHandler, err := container.Provide[*discord.DiscordHandler]()
	if err != nil {
		return err
	}

	if err := discordHandler.Open(ctx); err != nil {
		return err
	}

	return nil
}
