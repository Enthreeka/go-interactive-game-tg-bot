package view

import (
	"context"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/handler/tgbot"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/logger"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/tg/markup"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type ViewGeneral struct {
	log *logger.Logger
}

func NewViewGeneral(log *logger.Logger) *ViewGeneral {
	return &ViewGeneral{
		log: log,
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
		msg := tgbotapi.NewMessage(update.FromChat().ID, "Приветствуем на конкурсе… бла-бла (готовый текст предоставим), для участия нажмите на «Старт»")
		//msg.ReplyMarkup = &markup.StartMenu
		msg.ParseMode = tgbotapi.ModeHTML

		if _, err := bot.Send(msg); err != nil {
			v.log.Error("failed to send message", zap.Error(err))
			return err
		}

		return nil
	}
}
