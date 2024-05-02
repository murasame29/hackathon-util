package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func LoadEnv(filePath ...string) error {
	if err := godotenv.Load(filePath...); err != nil {
		return err
	}

	config := config{}

	if err := env.Parse(&config.Application); err != nil {
		return err
	}

	if err := env.Parse(&config.Spreadsheets); err != nil {
		return err
	}

	if err := env.Parse(&config.Discord); err != nil {
		return err
	}

	if err := env.Parse(&config.Google); err != nil {
		return err
	}

	Config = &config

	return nil
}
