package driver

import (
	"github.com/murasame29/hackathon-util/cmd/config"
	"google.golang.org/api/option"
)

func NewGCPCredential() option.ClientOption {
	return option.WithCredentialsFile(config.Config.Google.Credentials)
}
