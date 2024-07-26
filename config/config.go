package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken string
}

func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
func Load() *Config {
	dotenv := goDotEnvVariable("TELEGRAM_BOT_TOKEN")
	return &Config{
		BotToken: dotenv,
	}
}
