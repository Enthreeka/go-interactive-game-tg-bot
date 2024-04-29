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
		msg := tgbotapi.NewMessage(update.FromChat().ID, `Спасибо, теперь вы участник конкурса!
	
	Напоминаем условия:
	
	1. Конкурс продлится с 6 по 12 мая
	
	2. Вам нужно подписаться на каналы-участники и включить уведомления, чтобы не пропускать новые посты! Каждый пост нужно внимательно читать, чтобы получить подсказки для ответа на вопросы!
	
	Ещё раз эти каналы:
	
	Смарт инвестиции - https://t.me/+boSZNTXjmosxMGU6
	Фонда – https://t.me/+3HolQmLEA8k3ZGMy
	Простые инвестиции − https://t.me/+plO6P6Ye6fUwMjA6
	Биржа, деньги, хомяки − https://t.me/+sxWKWdMzElExZTZi
	Сигналы от души − https://t.me/+-q_q3w_j5DY1ZjAy
	Invest Assistance - https://t.me/+BAPMMRhFmckxNTk6
	
	3. Включите уведомления в самом боте, чтобы не пропустить новый вопрос! Вопросы мы будем задавать дважды в сутки на протяжении всего срока действия конкурса, а в последний день зададим сразу 3 вопроса!
	
	4. За каждый правильный ответ вы получите определённое количество баллов – чем сложнее вопрос, тем больше баллов!
	
	5. В конце конкурсной недели мы подведём итоги. Участники, заработавшие больше всех баллов, получат от нас призы:
	
	1 место – акция Полюса
	2-5 место – акция Лукойла
	6-10 место – акция Яндекса
	
	А также среди всех участников, набравших 5 баллов и более, будет проведён розыгрыш случайным образом – 10 победителей получат по 1 акции Татнефти!
	
	Важно: призы можно получить только акциями. Мы вышлем вам их на открытый счёт в Тинькофф. Если у вас там пока нет счёта, то откройте его по ссылке (https://partners.tinkoff.ru/click/87245a4a-e98b-43cd-b6bc-7255b91b9dc5). Вы получите 3 месяца бесплатного обслуживания и ещё бонусные акции от самого Тинькофф. Классно же!
	
	6. Если победителей, набравших равное количество баллов, будет несколько, мы проведём для них суперигру, чтобы выявить самых-самых!`)
		//msg.ReplyMarkup = &markup.StartMenu
		msg.ParseMode = tgbotapi.ModeHTML

		if _, err := bot.Send(msg); err != nil {
			v.log.Error("failed to send message", zap.Error(err))
			return err
		}

		return nil
	}
}
