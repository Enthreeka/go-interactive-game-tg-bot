package view

import (
	"context"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/handler/tgbot"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/service"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/logger"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/tg/markup"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type ViewGeneral struct {
	log      *logger.Logger
	commServ service.CommunicationService
}

func NewViewGeneral(log *logger.Logger, commServ service.CommunicationService) *ViewGeneral {
	return &ViewGeneral{
		log:      log,
		commServ: commServ,
	}
}

func (c *ViewGeneral) CallbackStartAdminPanel() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg := tgbotapi.NewMessage(update.FromChat().ID, "Панель управления")
		msg.ReplyMarkup = &markup.StartMenu
		msg.ParseMode = tgbotapi.ModeHTML

		if _, err := bot.Send(msg); err != nil {
			c.log.Error("failed to send message", zap.Error(err))
			return err
		}

		return nil
	}
}

func (v *ViewGeneral) ViewFirstMessage() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		commMsg, err := v.commServ.GetMessage(ctx)
		if err != nil {
			v.log.Error("failed to get communication msg: %v", err)
			return nil
		}

		msg := tgbotapi.NewMessage(update.FromChat().ID, commMsg)

		//msg.ReplyMarkup = &markup.StartMenu
		msg.ParseMode = tgbotapi.ModeHTML
		msg.DisableWebPagePreview = true

		if _, err := bot.Send(msg); err != nil {
			v.log.Error("failed to send message", zap.Error(err))
			return err
		}

		return nil
	}
}
