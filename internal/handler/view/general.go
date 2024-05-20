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
		msg := tgbotapi.NewMessage(update.FromChat().ID, `Друзья! Приветствуем вас на конкурсе «В преддверии лета»! 

В качестве призов разыграем акции. Призовой фонд такой:

1 место – акция Яндекса
2-5 место – акция Транснефти 
6-10 место – акция Инарктики 
11-15 место – акция Газпромнефти 

Среди 10 случайных участников, набравших 5 и более баллов, традиционно разыграем акции Татнефти! 

Акции будут перечислены на счёт Тинькофф. Если у вас не открыт там счёт, то вы можете сделать это по ссылке: https://partners.tinkoff.ru/click/87245a4a-e98b-43cd-b6bc-7255b91b9dc5

Условия:

1. Подписаться на каналы-спонсоры

Vyacheslav Goodwin https://t.me/+yu5PMjq4xwJkZGJi
Инвестор Альфа https://t.me/+sAdBODkhpQNmYjU6
Кравцова и рынки https://t.me/+fJmgyAaZKagzZjdi
Invest Premium https://t.me/+Qnlwt3uDpRQxY2My
Сигналы от души https://t.me/+zAnSVfHAShI0NGQy
Марафон инвестиций https://t.me/+gbZewNzODus2MGYy

2. Ждать вопросы по содержанию каналов – чем больше баллов наберёте, тем выше шанс на выигрыш (теперь вы понимаете, почему важно подписаться на все каналы?)

27 мая стартуем! И пусть победит сильнейший!

Отзывы о прошлом конкурсе и дополнительную информацию вы можете посмотреть по ссылке: https://t.me/tgresearh_contest.`)
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
