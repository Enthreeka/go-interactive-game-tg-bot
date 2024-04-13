package callback

import (
	"context"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/handler/tgbot"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/logger"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/store"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/tg"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/tg/markup"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackGeneral struct {
	log   *logger.Logger
	store *store.Store
	tgMsg *tg.TelegramMsg
}

func NewCallbackGeneral(log *logger.Logger, store *store.Store, tgMsg *tg.TelegramMsg) *CallbackGeneral {
	return &CallbackGeneral{
		log:   log,
		store: store,
		tgMsg: tgMsg,
	}
}

func (c *CallbackGeneral) CallbackCancelCommand() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		defer c.store.Delete(update.FromChat().ID)

		if err := c.tgMsg.SendNewMessage(update.FromChat().ID,
			nil,
			"Команда отменена",
		); err != nil {
			return err
		}

		return nil
	}
}

func (c *CallbackGeneral) CallbackMainMenu() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {

		if _, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&markup.StartMenu,
			"Панель управления",
		); err != nil {
			return err
		}

		return nil
	}
}
