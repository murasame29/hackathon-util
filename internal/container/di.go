package container

import (
	"net/http"

	"github.com/murasame29/hackathon-util/internal/application"
	"github.com/murasame29/hackathon-util/internal/driver"
	"github.com/murasame29/hackathon-util/internal/gateways/discordgo"
	"github.com/murasame29/hackathon-util/internal/gateways/gs"
	"github.com/murasame29/hackathon-util/internal/handler"
	"github.com/murasame29/hackathon-util/internal/router"
	"go.uber.org/dig"
)

var Container *dig.Container

type provideArg struct {
	constructor any
	opts        []dig.ProvideOption
}

func NewContainer() http.Handler {
	Container = dig.New()

	args := []provideArg{
		{driver.NewGCPCredential, []dig.ProvideOption{}},
		{discordgo.New, []dig.ProvideOption{}},
		{gs.New, []dig.ProvideOption{}},
		{router.NewRoute, []dig.ProvideOption{}},
		{handler.NewHandler, []dig.ProvideOption{}},
		{application.NewApplicationService, []dig.ProvideOption{}},
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
