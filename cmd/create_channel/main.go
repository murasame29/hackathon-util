package main

import (
	"context"
	"flag"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/murasame29/hackathon-util/cmd/config"
	"github.com/murasame29/hackathon-util/internal/gateways/discordgo"
	"github.com/murasame29/hackathon-util/internal/gateways/gs"
	"github.com/murasame29/hackathon-util/pkg/logger"
	"github.com/sourcegraph/conc"
	"google.golang.org/api/option"
)

var (
	credentialPath string
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
	flag.StringVar(&credentialPath, "c", "./credential.json", "path to credential.json file")
	flag.Parse()

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
	credential := option.WithCredentialsFile(credentialPath)

	ctx := logger.NewLoggerWithContext(context.Background())

	gs, err := gs.New(credential)
	if err != nil {
		logger.Error(ctx, "failed to create gs client", logger.Field("err", err))
		return err
	}

	values, err := gs.Read(config.Config.Spreadsheets.ID, config.Config.Spreadsheets.Range)
	if err != nil {
		logger.Error(ctx, "failed to read spreadsheet", logger.Field("err", err))
		return err
	}

	dg := discordgo.New()

	var (
		wg conc.WaitGroup
	)
	defer wg.Wait()

	for _, value := range values {
		teamName := value[0]

		// 正規表現を作成
		reg := regexp.MustCompile(`\(([^)]+)\)`)

		// 正規表現に一致する部分を空文字に置き換える
		name := reg.ReplaceAllString(teamName, "")

		if name == "欠席" {
			continue
		}

		logger.Info(ctx, "creating channel ", logger.Field("role", name))

		// 雑談 , 会議 , vcのチャンネルを作る

		wg.Go(func() {
			logger.Info(ctx, "creating category", logger.Field("category", name))
			categoryID, err := dg.CreateChannelCategory(ctx, config.Config.Discord.GuildID, name)
			if err != nil {
				logger.Error(ctx, "failed to create Category", logger.Field("err", err))
			}

			logger.Info(ctx, "Create category successful", logger.Field("category", name))
			logger.Info(ctx, "creating channel", logger.Field("channel", "雑談"))

			if _, err := dg.CreateChannelText(ctx, config.Config.Discord.GuildID, categoryID, "雑談"); err != nil {
				logger.Error(ctx, "failed to create channel", logger.Field("err", err))
			}

			logger.Info(ctx, "Create channel successful", logger.Field("channel", "雑談"))
			logger.Info(ctx, "creating channel", logger.Field("channel", "会議"))

			if _, err := dg.CreateChannelText(ctx, config.Config.Discord.GuildID, categoryID, "会議"); err != nil {
				logger.Error(ctx, "failed to create channel", logger.Field("err", err))
			}

			logger.Info(ctx, "Create channel successful", logger.Field("channel", "会議"))
			logger.Info(ctx, "creating channel", logger.Field("channel", "vc"))

			if _, err := dg.CreateChannelVoice(ctx, config.Config.Discord.GuildID, categoryID, "vc"); err != nil {
				logger.Error(ctx, "failed to create channel", logger.Field("err", err))
			}

			logger.Info(ctx, "Create channel successful", logger.Field("channel", "vc"))
		})
	}
	return nil
}
