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
		msg := tgbotapi.NewMessage(update.FromChat().ID, `Легендарный конкурс от агентства TG Research возвращается, и на этот раз с ещё более щедрыми призами. Что мы предлагаем?

Конкурс, на котором вам нужно отвечать на вопросы по содержанию каналов. Всё просто: правильный ответ = 1 балл. Чем больше баллов – тем лучше приз.

На этот раз мы разыграем:

1 место – акция Фосагро 
2-5 место – акция Самолёта 
6-10 место – акция Транснефти 
11-15 место – акция Новатэка 
16-20 место – акция Инарктики

Если победителей будет несколько, мы разыграем призовые акции через генератор случайных чисел среди тех, кто набрал максимум баллов!

А среди 15 участников, набравших 10 и более баллов, традиционно случайным образом разыграем акции Татнефти!

Кроме того, вас ждут подарки от наших спонсоров: доступы в вип-каналы и инвестиционные клубы, обучающие материалы, курсы. Ну классно же!

Напомним, что для участия нужно:

1. Подписаться на каналы-спонсоры:

Инвестор Альфа − https://t.me/+zOjxi4gpTuhjYjJi
Сигналы от души − https://t.me/+3ez1G0Zd2OQzODE6
Биржа, деньги, хомяки − https://t.me/+vGli79nqJUxmYmUy
Romanov Capital − https://t.me/+_O2M92JmCCdiZjBi
Инвестируй или проиграешь − https://t.me/+R-MmnXOUrfJiYjAy
Invest Whales − https://t.me/+mBVQ6ZxQeHY5NDli

2. Подписаться на официальную группу конкурса − https://t.me/tgresearh_contest (там же можно почитать отзывы о предыдущих розыгрышах)

Всё. Далее читаем каналы и ждём начала розыгрыша! Вопросы будут приходить непосредственно в данный бот.

Сам конкурс стартует 8 июля и продлится целую неделю – до 14 июля! 

Ждём всех и каждого!`)

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
