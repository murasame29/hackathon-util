package main

import (
	"context"
	"flag"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/murasame29/hackathon-util/cmd/config"
	"github.com/murasame29/hackathon-util/internal/discordgo"
	"github.com/murasame29/hackathon-util/internal/gs"
	"github.com/murasame29/hackathon-util/internal/pkg/logger"
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

	var wg conc.WaitGroup
	defer wg.Wait()

	for _, value := range values {
		teamName := value[0]

		// 正規表現を作成
		reg := regexp.MustCompile(`\(([^)]+)\)`)

		// 正規表現に一致する部分を空文字に置き換える
		role := reg.ReplaceAllString(teamName, "")

		if role == "欠席" {
			continue
		}

		logger.Info(ctx, "creating role ", logger.Field("role", role))

		wg.Go(func() {
			if err := dg.CreateRole(ctx, role); err != nil {
				logger.Error(ctx, "failed to create Role", logger.Field("role", role))
			}

			logger.Info(ctx, "Role Create Successful", logger.Field("role", role))
		})
	}
	return nil
}
