package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	BotToken string
	Port     string
	ApiToken string
	Postgres Postgres
}

type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func LoadConfig(path string) Config {
	godotenv.Load(path + "/.env") // Load .env file if it exists

	conf := viper.New()
	conf.AutomaticEnv()

	cfg := Config{
		BotToken: conf.GetString("TG_BOT"),
		Port:     conf.GetString("HTTP_PORT"),
		ApiToken: conf.GetString("API_TOKEN"),
		Postgres: Postgres{
			Host:     conf.GetString("POSTGRES_HOST"),
			Port:     conf.GetString("POSTGRES_PORT"),
			User:     conf.GetString("POSTGRES_USER"),
			Password: conf.GetString("POSTGRES_PASSWORD"),
			Database: conf.GetString("POSTGRES_DATABASE"),
		},
	}

	return cfg
}
