package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/murasame29/hackathon-util/cmd/config"
	"github.com/murasame29/hackathon-util/internal/container"
	"github.com/murasame29/hackathon-util/internal/server"
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

	// show all dir
	log.Println(os.Getwd())
	search()

	handler := container.NewContainer()

	server.
		New(config.Config.Application.Addr, handler).
		RunWithGraceful(ctx)

	return nil
}

func search() {
	dir, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range dir {
		if d.IsDir() {
			search()
		}
		log.Println("f:", d.Name())
	}
}
