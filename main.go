package main

import (
	"log"

	"github.com/BaboaAtCity/BigNewsMorgan/bot"
	"github.com/BaboaAtCity/BigNewsMorgan/config"
)

func main() {
	cfg := config.Load()
	b, err := bot.New(cfg.BotToken)
	if err != nil {
		log.Panic(err)
	}

	b.Start()
}
