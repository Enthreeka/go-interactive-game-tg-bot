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
		msg := tgbotapi.NewMessage(update.FromChat().ID, `Рады приветствовать вас на третьем сезоне розыгрыша TG Research! 

			В качестве призов разыграем, как всегда, акции. На этот раз призовой фонд такой:
			
			1 место – акция Самолёта
			2-5 место – акция Тинькофф
			6-10 место – акция Северстали
			11-15 место – акция Новатэк
			
			Если победителей будет несколько, то мы разыграем призовые акции через генератор случайных чисел среди тех, кто набрал максимум очков!
			
			А также среди 10 случайных участников, набравших 10 и более баллов, мы традиционно разыграем акции Татнефти! 
			
			Условия:
			
			1. Подписаться на каналы-спонсоры
			
			Invest Smart – https://t.me/+FlDh2v_TCMEyNmUy 
			INSpace – https://t.me/+AsHtAzik6poxMWIx 
			Frolov I&R – https://t.me/+9vh20vIAIrc5MjE0 
			Первый инвестиционный – https://t.me/+jKIA54kMgYY5MmQy  
			Капибара на бирже – https://t.me/+vM_MxX_6bZViNmJi 
			Unique Trade – https://t.me/+fW80I16SdbwwNzky 
			Приватка Владимир на бирже − https://t.me/+WrbxtN_P8mxmYmQ8 
			
			2. Ждать вопросы – их зададут прямо в этом боте.
			
			Конкурс продлится с 10 по 16 июня. В последний день в 20:00 мск подведём итоги. Призы раздадим 17-20 июня.
			
			Всем искренне желаем удачи!
			
			Подробности в официальной группе розыгрыша: https://t.me/tgresearh_contest.`)

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
