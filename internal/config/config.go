package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

// validate is a package-level singleton to reuse cached struct metadata.
var validate = validator.New()

// Config is the top-level structure that maps to the YAML manifest.
//
// Example:
//
//	eventName: "hoge"
//	googleSheet:
//	  id: "hoge"
//	  teamTableRange: "hoge!A:Z"
//	  credentialFile: "./credential.json"
//	discord:
//	  guildID: "hoge"
//	  muteRoleID: "hoge"
//	  enablePrivateVC: false
//	  enablePrivateCategory: false
type Config struct {
	EventName   string        `yaml:"eventName"   validate:"required"`
	GoogleSheet SheetConfig   `yaml:"googleSheet" validate:"required"`
	Discord     DiscordConfig `yaml:"discord"     validate:"required"`
}

// SheetConfig holds Google Sheets related settings.
type SheetConfig struct {
	ID             string `yaml:"id"             validate:"required"`
	TeamTableRange string `yaml:"teamTableRange" validate:"required"`
	CredentialFile string `yaml:"credentialFile" validate:"required"`
}

// DiscordConfig holds Discord related settings.
type DiscordConfig struct {
	GuildID               string `yaml:"guildID"               validate:"required"`
	MuteRoleID            string `yaml:"muteRoleID"`
	EnablePrivateVC       bool   `yaml:"enablePrivateVC"`
	EnablePrivateCategory bool   `yaml:"enablePrivateCategory"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %s: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	dec := yaml.NewDecoder(f)
	dec.KnownFields(true) // reject unknown keys
	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: decode %s: %w", path, err)
	}

	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("config: validation failed: %w", err)
	}

	return &cfg, nil
}
