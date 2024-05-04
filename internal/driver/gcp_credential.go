package driver

import (
	"log"
	"os"

	"github.com/murasame29/hackathon-util/cmd/config"
	"google.golang.org/api/option"
)

func NewGCPCredential() option.ClientOption {
	if config.Config.Google.Credentials == "" {
		return nil
	}

	if _, err := os.Open(config.Config.Google.Credentials); err != nil {
		de, err := os.ReadDir(".")
		if err != nil {
			log.Println(err)
		}

		log.Println("credentials file not found")
		log.Println("de;", de)

		for _, d := range de {
			log.Println(d.Name())
		}

		log.Println(err)
	}
	return option.WithCredentialsFile(config.Config.Google.Credentials)
}
