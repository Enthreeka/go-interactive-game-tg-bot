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
		msg := tgbotapi.NewMessage(update.FromChat().ID, `Рады приветствовать вас на 5 сезоне конкурса TG Research! Напомним, что для победы нужно отвечать на вопросы прямо в этом боте.

Победители получат:

1-3 место - акция Черкизово 
4-10 место - акция Северстали 
10-15 место - акция Новатэка 
16-20 место - акция Инарктики 

А среди 20 случайных участников, набравших 10 баллов и более, мы традиционно разыграем акции Татнефти!

Также победители и участники получат доступы в вип-каналы, бесплатные образовательные курсы, обучающую литературу и множество других призов!

Для участия нужно подписаться на каналы-спонсоры: 

Инвест Эра - https://t.me/+z--4_BPTPjdlYWYy 
Ракета инвестора - https://t.me/+q2a9vT5fpG85YTMy 
Биржевой маклер - https://t.me/+EeKodfU2BqJjNmFi
Кравцова и рынки - https://t.me/+3d8Ok49qx-FmMTky 
Биржевая ключница - https://t.me/+osJ27uNmnUg5ZjMy 

С 1 августа вы будете получать вопросы. Итоги подведём 11 августа, в воскресенье днём.

Отзывы о предыдущих сезонах и больше подробностей - в официальной группе конкурса https://t.me/tgresearh_contest.

Увидимся! Ваша команда TG Research.`)

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
