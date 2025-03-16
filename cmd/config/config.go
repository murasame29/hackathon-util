package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func NewEnvironment() *EnvironmentsVariables {
	return &EnvironmentsVariables{}
}

func LoadEnv(filePath ...string) error {
	if len(filePath) > 0 {
		if err := godotenv.Load(filePath...); err != nil {
			return err
		}
	}

	config := EnvironmentsVariables{}

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
