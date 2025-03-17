package container

import (
	"github.com/murasame29/hackathon-util/internal/driver"
	"github.com/murasame29/hackathon-util/internal/framewrok/discord"
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
		{discord.NewHandler, []dig.ProvideOption{}},
	}

	for _, arg := range args {
		if err := Container.Provide(arg.constructor, arg.opts...); err != nil {
			return err
		}
	}

	return nil
}

func Provide[T any]() (T, error) {
	var t T
	if err := Container.Invoke(func(v T) {
		t = v
	}); err != nil {
		return t, err
	}
	return t, nil
}
