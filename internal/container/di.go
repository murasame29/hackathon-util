package container

import (
	"net/http"

	"github.com/murasame29/hackathon-util/internal/adapter/controller"
	"github.com/murasame29/hackathon-util/internal/adapter/gateways/discordgo"
	"github.com/murasame29/hackathon-util/internal/adapter/gateways/gs"
	"github.com/murasame29/hackathon-util/internal/application"
	"github.com/murasame29/hackathon-util/internal/driver"
	"github.com/murasame29/hackathon-util/internal/framewrok/discord"
	"github.com/murasame29/hackathon-util/internal/framewrok/http/router"
	"go.uber.org/dig"
)

var Container *dig.Container

type provideArg struct {
	constructor any
	opts        []dig.ProvideOption
}

func NewContainer() error {
	Container = dig.New()

	args := []provideArg{
		{driver.NewDiscordSession, []dig.ProvideOption{}},
		{driver.NewGCPCredential, []dig.ProvideOption{}},
		{discordgo.New, []dig.ProvideOption{}},
		{gs.New, []dig.ProvideOption{}},
		{router.NewRoute, []dig.ProvideOption{}},
		{discord.NewHandler, []dig.ProvideOption{}},
		{controller.NewHandler, []dig.ProvideOption{}},
		{application.NewApplicationService, []dig.ProvideOption{}},
	}

	for _, arg := range args {
		if err := Container.Provide(arg.constructor, arg.opts...); err != nil {
			return err
		}
	}

	return nil
}

func NewSheetLessContainer() http.Handler {
	Container = dig.New()

	args := []provideArg{
		{driver.NewDiscordSession, []dig.ProvideOption{}},
		{discordgo.New, []dig.ProvideOption{}},
		{router.NewRoute, []dig.ProvideOption{}},
		{discord.NewHandler, []dig.ProvideOption{}},
		{controller.NewHandler, []dig.ProvideOption{}},
		{application.NewSheetLessApplicationService, []dig.ProvideOption{}},
	}

	for _, arg := range args {
		if err := Container.Provide(arg.constructor, arg.opts...); err != nil {
			panic(err)
		}
	}

	var handler http.Handler
	if err := Container.Invoke(func(h http.Handler) {
		handler = h
	}); err != nil {
		panic(err)
	}

	return handler
}

func Provide(fn any) error {
	return Container.Invoke(fn)
}
