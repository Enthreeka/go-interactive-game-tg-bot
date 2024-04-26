package callback

import (
	"context"
	"fmt"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/boterror"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/entity"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/handler"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/handler/tgbot"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/service"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/excel"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/logger"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/store"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/tg"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/tg/markup"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"sync"
)

type CallbackContest struct {
	contestService service.ContestService
	userService    service.UserService
	store          *store.Store
	log            *logger.Logger
	tgMsg          *tg.TelegramMsg
	excel          *excel.Excel

	mu sync.RWMutex
}

func NewCallbackContest(contestService service.ContestService, userService service.UserService, store *store.Store, log *logger.Logger, tgMsg *tg.TelegramMsg, excel *excel.Excel) *CallbackContest {
	return &CallbackContest{
		contestService: contestService,
		userService:    userService,
		log:            log,
		tgMsg:          tgMsg,
		store:          store,
		excel:          excel,
	}
}

// CallbackContestSetting - contest_setting
func (c *CallbackContest) CallbackContestSetting() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {

		if _, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&markup.ContestSetting,
			"CallbackContestSetting"); err != nil {
			return err
		}

		return nil
	}
}

// CallbackGetAllContest - get_all_contest
// Выводятся кнопки вида contest_get_{contest_id}
func (c *CallbackContest) CallbackGetAllContest() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		_, markupContest, err := c.contestService.GetAllContestsButtons(ctx, "get")
		if err != nil {
			c.log.Error("contestService.GetAllContests: failed to get contest: %v", err)
			return err
		}

		if _, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			markupContest,
			"CallbackGetAllContest"); err != nil {
			return err
		}

		return nil
	}
}

// CallbackCreateContest - create_contest
func (c *CallbackContest) CallbackCreateContest() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&markup.CancelState,
			`Отправить через json`+
				"\n{\n"+
				`"название_конкурса":"сюда вписать название",`+"\n"+
				`"дедлайн":"2006-01-02 15:04:05"`+
				"\n}")
		if err != nil {
			return err
		}

		c.store.Delete(update.FromChat().ID)
		c.store.Set(&store.ContestStore{
			MsgID:              msg,
			UserID:             update.FromChat().ID,
			TypeCommandContest: store.ContestCreate,
		}, update.FromChat().ID)

		return nil
	}
}

// CallbackDeleteContest - delete_contest
func (c *CallbackContest) CallbackDeleteContest() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		_, markupContest, err := c.contestService.GetAllContestsButtons(ctx, "delete")
		if err != nil {
			c.log.Error("contestService.GetAllContests: failed to get contest: %v", err)
			return err
		}

		if _, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			markupContest,
			"CallbackDeleteContest"); err != nil {
			return err
		}

		return nil
	}
}

// CallbackCreateMailing - create_mailing
func (c *CallbackContest) CallbackCreateMailing() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		msg, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&markup.CancelState,
			"Для моментальной рассылки укажите текст. Отправьте следующее сообщение в виде json.\n{\n  \"сообщение\": \"сюда нужно вписать текст\"\n}")
		if err != nil {
			return err
		}

		c.store.Delete(update.FromChat().ID)
		c.store.Set(&store.ContestStore{
			MsgID:              msg,
			UserID:             update.FromChat().ID,
			TypeCommandContest: store.CreateUserMailing,
		}, update.FromChat().ID)

		return nil
	}
}

// CallbackGetContestByID - contest_get_{contest_id}
func (c *CallbackContest) CallbackGetContestByID() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		contestID := entity.GetContestID(update.CallbackData())

		contest, err := c.contestService.GetContestByID(ctx, contestID)
		if err != nil {
			c.log.Error("contestService.GetContestByID: failed to get contest: %v", err)
			return err
		}

		markupContest := markup.ContestByIDSetting(contestID)
		if _, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&markupContest,
			fmt.Sprintf("%v:", contest),
		); err != nil {
			return err
		}

		return nil
	}
}

// CallbackDownloadRating - download_rating_{contest_id}
func (c *CallbackContest) CallbackDownloadRating() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		contestID := entity.GetContestID(update.CallbackData())

		userResult, err := c.contestService.GetAllUserResultsByContest(ctx, contestID)
		if err != nil {
			c.log.Error("contestService.GetAllUserResultsByContest: failed to get contest: %v", err)
			return err
		}

		c.mu.Lock()
		fileName, err := c.excel.GenerateUserResultsExcelFile(userResult, contestID, update.CallbackQuery.From.UserName)
		if err != nil {
			c.log.Error("Excel.GenerateExcelFile: failed to generate excel file: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		fileIDBytes, err := c.excel.GetExcelFile(fileName)
		if err != nil {
			c.log.Error("Excel.GetExcelFile: failed to get excel file: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}
		c.mu.Unlock()

		if fileIDBytes == nil {
			c.log.Error("fileIDBytes: %v", boterror.ErrNil)
			handler.HandleError(bot, update, boterror.ParseErrToText(boterror.ErrNil))
			return nil
		}

		if _, err := c.tgMsg.SendDocument(update.FromChat().ID,
			fileName,
			fileIDBytes,
			"CallbackDownloadRating",
		); err != nil {
			c.log.Error("tgMsg.SendDocument: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		return nil
	}
}

// CallbackContestDelete - contest_delete_{contest_id}
func (c *CallbackContest) CallbackContestDelete() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		contestID := entity.GetContestID(update.CallbackData())

		if err := c.contestService.DeleteContest(ctx, contestID); err != nil {
			c.log.Error("contestService.DeleteContest: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		if _, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			nil,
			"Удалено успешно",
		); err != nil {
			return err
		}

		return nil
	}
}

// CallbackContestReminder - contest_reminder_{contest_id}
func (c *CallbackContest) CallbackContestReminder() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		contestID := entity.GetContestID(update.CallbackData())

		contest, err := c.contestService.GetContestByID(ctx, contestID)
		if err != nil {
			c.log.Error("contestService.GetContestByID: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		users, err := c.userService.GetAllUsers(ctx)
		if err != nil {
			c.log.Error("userService.GetAllUsers: failed to get users: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		go func(u []entity.User, adminID int64) {
			h, m, _ := contest.Deadline.Clock()

			var minu string
			if m > 9 {
				minu = strconv.Itoa(m)
			} else {
				minu = strconv.Itoa(m)
				minu = "0" + minu
			}
			text := fmt.Sprintf("Дорогие уастники! Напоминаем Вам, что у нас проходит конкурс: %s.\n"+
				"Он завершится: %d числа в %d:%s.", contest.Name, contest.Deadline.Day(), h, minu)

			var totalSend int

			for _, user := range u {

				if user.BlockedBot == false {
					if err := c.tgMsg.SendNewMessage(user.ID, nil, text); err != nil {
						c.log.Error("tgMsg.SendNewMessage in user question send: %v", err)

						if strings.Contains(err.Error(), "Forbidden: bot was blocked by the user") ||
							strings.Contains(err.Error(), "Bad Request: chat not found") {

							if err := c.userService.UpdateBlockedBotStatus(context.Background(), user.ID, true); err != nil {
								c.log.Error("userService.UpdateBlockedBotStatus: %v", err)
							}

						} else {
							c.log.Error("error on sending: %v", err)
						}
					}
					totalSend++
				}

			}

			if err := c.tgMsg.SendNewMessage(
				adminID,
				nil,
				fmt.Sprintf("Рассылка завершена. Отправлено пользователям: %d", totalSend),
			); err != nil {
				c.log.Error("questionsService.UpdateIsSendByQuestionID: %v", err)
				return
			}
		}(users, update.FromChat().ID)

		return nil
	}
}

// CallbackSendRating - send_rating_{contest_id}
// Отправляется только людям, которые участвовали в текущем контексте
func (c *CallbackContest) CallbackSendRating() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		contestID := entity.GetContestID(update.CallbackData())

		userResult, err := c.contestService.GetAllUserResultsByContest(ctx, contestID)
		if err != nil {
			c.log.Error("contestService.GetAllUserResultsByContest: failed to get contest: %v", err)
			return err
		}

		c.mu.Lock()
		fileName, err := c.excel.GenerateForUserResultsExcelFile(userResult, contestID, update.CallbackQuery.From.UserName)
		if err != nil {
			c.log.Error("Excel.GenerateExcelFile: failed to generate excel file: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		fileIDBytes, err := c.excel.GetExcelFile(fileName)
		if err != nil {
			c.log.Error("Excel.GetExcelFile: failed to get excel file: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}
		c.mu.Unlock()

		go func(u []entity.UserResult, adminID int64) {
			var totalSend int

			for _, user := range u {

				if _, err := c.tgMsg.SendDocument(user.UserID, fileName, fileIDBytes,
					"Отправляем Вам таблицу с результатами пользователей за текущий конкурс"); err != nil {
					c.log.Error("tgMsg.SendDocument to user table: %v", err)

					if strings.Contains(err.Error(), "Forbidden: bot was blocked by the user") ||
						strings.Contains(err.Error(), "Bad Request: chat not found") {

						if err := c.userService.UpdateBlockedBotStatus(context.Background(), user.UserID, true); err != nil {
							c.log.Error("userService.UpdateBlockedBotStatus: %v", err)
						}

					} else {
						c.log.Error("error on sending: %v", err)
					}
				}
				totalSend++
			}

			if err := c.tgMsg.SendNewMessage(
				adminID,
				nil,
				fmt.Sprintf("Рассылка завершена. Отправлено пользователям: %d", totalSend),
			); err != nil {
				c.log.Error("questionsService.UpdateIsSendByQuestionID: %v", err)
				return
			}
		}(userResult, update.FromChat().ID)

		return nil
	}
}

// CallbackPickRandom - pick_random_{contest_id}
func (c *CallbackContest) CallbackPickRandom() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		contestID := entity.GetContestID(update.CallbackData())

		msg, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&markup.CancelState,
			`Отправьте сообщение в следующем формате:`+
				"\n{\n"+
				`"рейтинг": сюда нужно вписать целое число,`+"\n"+
				`"количество_людей": сюда нужно вписать целое число`+
				"\n}")
		if err != nil {
			return err
		}

		c.store.Delete(update.FromChat().ID)
		c.store.Set(&store.ContestStore{
			MsgID:              msg,
			UserID:             update.FromChat().ID,
			ContestID:          contestID,
			TypeCommandContest: store.ContestPick,
		}, update.FromChat().ID)

		return nil
	}
}

// CallbackSendMessage - send_message_{contest_id}
// Диалоговое окно с пользваотелем
func (c *CallbackContest) CallbackSendMessage() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		contestID := entity.GetContestID(update.CallbackData())

		msg, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&markup.CancelState,
			`Начните диалоговое окно с пользователем. Пользователь не сможет отправить вам ответ в бота. Отправьте сообщение в следующем формате:`+
				"\n{\n"+
				`"сообщение":"впишите сюда текст",`+"\n"+
				`"id_пользователя": сюда нужно вписать целое число`+
				"\n}")
		if err != nil {
			return err
		}

		c.store.Delete(update.FromChat().ID)
		c.store.Set(&store.ContestStore{
			MsgID:              msg,
			UserID:             update.FromChat().ID,
			ContestID:          contestID,
			TypeCommandContest: store.ContestUser,
		}, update.FromChat().ID)

		return nil
	}
}

// CallbackUpdateRating - update_rating_{contest_id}
func (c *CallbackContest) CallbackUpdateRating() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		contestID := entity.GetContestID(update.CallbackData())

		msg, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&markup.CancelState,
			`Вы можете изменить рейтинг конкретному пользователю. Отправьте сообщение в следующем формате:`+
				"\n{\n"+
				`"рейтинг":сюда нужно вписать целое число,`+"\n"+
				`"id_пользователя": сюда нужно вписать целое число`+
				"\n}")
		if err != nil {
			return err
		}

		c.store.Delete(update.FromChat().ID)
		c.store.Set(&store.ContestStore{
			MsgID:              msg,
			UserID:             update.FromChat().ID,
			ContestID:          contestID,
			TypeCommandContest: store.ContestRating,
		}, update.FromChat().ID)

		return nil
	}
}
