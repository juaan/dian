package config

import (
	"os"

	"github.com/subosito/gotenv"
)

type (
	Config struct {
		TailvyToken  string
		DiscordToken string
	}
)

var (
	config *Config
)

func InitConfig() {
	gotenv.Load()

	config = &Config{
		TailvyToken:  os.Getenv("TAILVY_TOKEN"),
		DiscordToken: os.Getenv("DISCORD_TOKEN"),
	}
}

func Get() *Config {
	return config
}
