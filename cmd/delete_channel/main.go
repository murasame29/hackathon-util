package main

import (
	"context"
	"flag"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/murasame29/hackathon-util/cmd/config"
	"github.com/murasame29/hackathon-util/internal/discordgo"
	"github.com/murasame29/hackathon-util/internal/gs"
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

	dg, err := discordgo.New()
	if err != nil {
		logger.Error(ctx, "failed to create discord client", logger.Field("err", err))
		return err
	}

	channels, err := dg.GetChannel(ctx)
	if err != nil {
		logger.Error(ctx, "failed to get channels", logger.Field("err", err))
		return err
	}

	roles := make(map[string]struct{})
	for _, value := range values {
		teamName := value[0]

		// 正規表現を作成
		reg := regexp.MustCompile(`\(([^)]+)\)`)

		// 正規表現に一致する部分を空文字に置き換える
		role := reg.ReplaceAllString(teamName, "")

		if role == "欠席" {
			continue
		}

		roles[strings.TrimSpace(role)] = struct{}{}
	}

	var deleteTargetChannelIDs []string
	currentCategory := make(map[string]string)

	for _, channel := range channels {
		if channel.Type == 4 {
			// rolesに存在しない場合無視
			if _, ok := roles[channel.Name]; !ok {
				continue
			}

			currentCategory[channel.ID] = channel.Name
			deleteTargetChannelIDs = append(deleteTargetChannelIDs, channel.ID)
			continue
		}

		_, ok := currentCategory[channel.ParentID]
		if !ok {
			continue
		}

		deleteTargetChannelIDs = append(deleteTargetChannelIDs, channel.ID)
	}

	var (
		wg conc.WaitGroup
	)
	defer wg.Wait()

	for _, deleteTargetChannelID := range deleteTargetChannelIDs {
		wg.Go(func() {
			time.Sleep(time.Second / 10)
			logger.Info(ctx, "deleting channel ", logger.Field("channelId", deleteTargetChannelID))

			if err := dg.DeleteChannel(ctx, deleteTargetChannelID); err != nil {
				logger.Error(ctx, "failed to delete channel", logger.Field("err", err))
			}

			logger.Info(ctx, "channel delete successful", logger.Field("channelId", deleteTargetChannelID))
		})
	}
	return nil
}
