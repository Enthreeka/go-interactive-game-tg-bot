package main

import (
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/bot"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/config"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/logger"
)

func main() {
	log := logger.New()

	cfg, err := config.New()
	if err != nil {
		log.Fatal("failed load config: %v", err)
	}

	newBot := bot.NewBot()

	if err := newBot.Run(log, cfg); err != nil {
		log.Fatal("failed to run telegram bot: %v", err)
	}
}
