package bot

import (
	"context"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/config"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/handler/tgbot"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/logger"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/postgres"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"os/signal"
	"syscall"
)

type Bot struct {
	psql *postgres.Postgres
}

func NewBot() *Bot {
	return &Bot{}
}

func (b *Bot) initialize(ctx context.Context, log *logger.Logger) {

}

func (b *Bot) Run(log *logger.Logger, cfg *config.Config) error {
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Fatal("failed to load token %v", err)
	}
	bot.Debug = false

	log.Info("Authorized on account %s", bot.Self.UserName)

	psql, err := postgres.New(context.Background(), 5, cfg.Postgres.URL)
	if err != nil {
		log.Fatal("failed to connect PostgreSQL: %v", err)
	}
	defer psql.Close()
	b.psql = psql

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	b.initialize(ctx, log)

	newBot := tgbot.NewBot(bot, log)

	if err := newBot.Run(ctx); err != nil {
		log.Error("failed to run tgbot: %v", err)
	}
	return nil
}
