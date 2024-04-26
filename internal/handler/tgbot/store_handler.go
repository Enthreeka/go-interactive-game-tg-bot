package tgbot

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/entity"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/store"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/tg/markup"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	"math/rand"
	"strings"
	"time"
)

func ParseJSON[T any](src string) (T, error) {
	var args T
	if err := json.Unmarshal([]byte(src), &args); err != nil {
		return *(new(T)), err
	}
	return args, nil
}

func (b *Bot) isStateExist(userID int64) (interface{}, bool) {
	data, exist := b.store.Read(userID)
	return data, exist
}

func (b *Bot) isStoreProcessing(ctx context.Context, update *tgbotapi.Update) (bool, error) {
	userID := update.Message.From.ID
	storeData, isExist := b.isStateExist(userID)
	if !isExist {
		return false, nil
	}
	defer b.store.Delete(userID)

	switch data := storeData.(type) {
	case *store.ContestStore:
		b.log.Info("process store.ContestStore")

		if data.TypeCommandContest == store.ContestCreate {
			if err := b.contentStoreCreate(ctx, update); err != nil {
				b.log.Error("contentStoreCreate: %v", err)
				return true, err
			}

			if err := b.updateChatInSuccessfullyCase(userID, data.MsgID, update.Message.MessageID, markup.ContestSetting, "И тут тоже поменять"); err != nil {
				return true, err
			}
		}

		if data.TypeCommandContest == store.ContestPick {
			if err := b.pickUsers(ctx, update, data); err != nil {
				b.log.Error("pickUsers: %v", err)
				return true, err
			}
		}

		if data.TypeCommandContest == store.ContestUser {
			if err := b.sendMessageToUser(ctx, update, data); err != nil {
				b.log.Error("sendMessageToUser: %v", err)
				return true, err
			}
		}

		if data.TypeCommandContest == store.ContestRating {
			if err := b.updateRating(ctx, update, data); err != nil {
				b.log.Error("updateRating: %v", err)
				return true, err
			}
		}

		if data.TypeCommandContest == store.CreateUserMailing {
			if err := b.userMailing(ctx, update); err != nil {
				b.log.Error("userMailing: %v", err)
				return true, err
			}
		}

		return true, nil
	case *store.QuestionStore:
		b.log.Info("process store.QuestionStore")

		if err := b.switchTypeCommandQuestion(ctx, data, update); err != nil {
			b.log.Error("switchTypeCommandQuestion: failed to process store.QuestionStore: %v", err)
			return true, err
		}

		//if err := b.updateChatInSuccessfullyCase(userID, data.MsgID, update.Message.MessageID); err != nil {
		//	return true, err
		//}

		return true, nil
	case store.AdminStore:
		defer b.store.Delete(userID)

		if data.TypeCommand == store.UserAdminCreate {
			if err := b.userService.UpdateRoleByUsername(ctx, "admin", update.Message.Text); err != nil {
				b.log.Error("isStoreExist:.UpdateRoleByUsername: %v", err)
				return true, err
			}
		}
		if data.TypeCommand == store.UserAdminDelete {
			if err := b.userService.UpdateRoleByUsername(ctx, "user", update.Message.Text); err != nil {
				b.log.Error("isStoreExist:userRepo.UpdateRoleByUsername: %v", err)
				return true, err
			}
		}

		b.log.Error("isStoreExist: undefind type command")
		return true, nil
	default:
		b.log.Error("undefined type data in storeProcessing")
		return true, nil
	}
}

func (b *Bot) switchTypeCommandQuestion(ctx context.Context, data *store.QuestionStore, update *tgbotapi.Update) error {
	switch data.TypeCommandQuestion {
	case store.QuestionCreate:
		if err := b.questionsService.CreateQuestion(ctx, &entity.Question{
			ContestID:     data.ContestID,
			CreatedByUser: data.UserID,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			QuestionName:  update.Message.Text,
		}); err != nil {
			b.log.Error("questionsService.CreateQuestion: %v", err)
			return err
		}

		markupQuestion := markup.QuestionSetting(data.ContestID)
		if err := b.updateChatInSuccessfullyCase(
			update.FromChat().ID,
			data.MsgID,
			update.Message.MessageID,
			markupQuestion,
			"И тут тоже поменять",
		); err != nil {
			b.log.Error("updateChatInSuccessfullyCase: %v", err)
			return err
		}
	case store.QuestionUpdate:
		if err := b.questionsService.UpdateQuestionName(ctx, data.QuestionID, update.Message.Text); err != nil {
			b.log.Error("questionsService.UpdateQuestionName: %v", err)
			return err
		}

		contestID, err := b.answersService.GetContestIDByQuestionID(ctx, data.QuestionID)
		if err != nil {
			b.log.Error("answersService.GetContestIDByQuestionID: %v", err)
			return err
		}

		questionMarkup := markup.QuestionByIDSetting(data.QuestionID, contestID)
		if err := b.updateChatInSuccessfullyCase(
			update.FromChat().ID,
			data.MsgID,
			update.Message.MessageID,
			questionMarkup,
			"И тут тоже поменять",
		); err != nil {
			b.log.Error("updateChatInSuccessfullyCase: %v", err)
			return err
		}
	case store.QuestionAddButtonAnswer:
		if err := b.contentStoreAddAnswer(ctx, update, data.QuestionID); err != nil {
			b.log.Error("contentStoreAddAnswer: %v", err)
			return err
		}

		contestID, err := b.answersService.GetContestIDByQuestionID(ctx, data.QuestionID)
		if err != nil {
			b.log.Error("answersService.GetContestIDByQuestionID: %v", err)
			return err
		}

		questionMarkup := markup.QuestionByIDSetting(data.QuestionID, contestID)
		if err := b.updateChatInSuccessfullyCase(
			update.FromChat().ID,
			data.MsgID,
			update.Message.MessageID,
			questionMarkup,
			"И тут тоже поменять",
		); err != nil {
			b.log.Error("updateChatInSuccessfullyCase: %v", err)
			return err
		}
	case store.QuestionAddDeadline:
		parsedTime, err := time.Parse(time.DateTime, update.Message.Text)
		if err != nil {
			return err
		}

		if err := b.questionsService.UpdateDeadlineByQuestionID(ctx, data.QuestionID, parsedTime); err != nil {
			b.log.Error("questionsService.UpdateDeadlineByQuestionID: %v", err)
			return err
		}

		contestID, err := b.answersService.GetContestIDByQuestionID(ctx, data.QuestionID)
		if err != nil {
			b.log.Error("answersService.GetContestIDByQuestionID: %v", err)
			return err
		}

		questionMarkup := markup.QuestionByIDSetting(data.QuestionID, contestID)
		if err := b.updateChatInSuccessfullyCase(
			update.FromChat().ID,
			data.MsgID,
			update.Message.MessageID,
			questionMarkup,
			"И тут тоже поменять",
		); err != nil {
			b.log.Error("updateChatInSuccessfullyCase: %v", err)
			return err
		}

	case store.QuestionTop10:

		if err := b.additionQuestion(ctx, update, data); err != nil {
			b.log.Error("additionQuestion: %v", err)
			return err
		}

	}

	return nil
}

func (b *Bot) contentStoreAddAnswer(ctx context.Context, update *tgbotapi.Update, questionID int) error {
	args, err := ParseJSON[entity.ArgsAddButton](update.Message.Text)
	if err != nil {
		b.log.Error("ParseJSON: %v", err)
		return err
	}

	tx, err := b.pg.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	if err := b.answersService.CreateAnswer(ctx, tx, &entity.Answer{
		Answer:         args.NameAnswer,
		CostOfResponse: args.CostOfAnswer,
		QuestionID:     questionID,
	}); err != nil {
		b.log.Error("answersService.CreateAnswer: %v", err)
		return err
	}

	return nil
}

func (b *Bot) contentStoreCreate(ctx context.Context, update *tgbotapi.Update) error {
	args, err := ParseJSON[entity.ArgsCreate](update.Message.Text)
	if err != nil {
		b.log.Error("ParseJSON: %v", err)
		return err
	}

	parsedTime, err := time.Parse(time.DateTime, args.Deadline)
	if err != nil {
		return err
	}

	if err = b.contestService.CreateContest(ctx, &entity.Contest{
		Name:     args.Name,
		Deadline: parsedTime,
	}); err != nil {
		return err
	}

	return nil
}

func (b *Bot) updateChatInSuccessfullyCase(userID int64, telegramMessageID, userMessageID int, markup tgbotapi.InlineKeyboardMarkup, text string) error {
	if resp, err := b.bot.Request(tgbotapi.NewDeleteMessage(userID,
		userMessageID)); nil != err || !resp.Ok {
		b.log.Error("failed to delete message id %d (%s): %v", userMessageID, string(resp.Result), err)
		return err
	}

	if _, err := b.tgMsg.SendEditMessage(userID, telegramMessageID, &markup, text); err != nil {
		return err
	}

	return nil
}

func (b *Bot) pickUsers(ctx context.Context, update *tgbotapi.Update, data *store.ContestStore) error {
	args, err := ParseJSON[entity.ArgsPick](update.Message.Text)
	if err != nil {
		b.log.Error("ParseJSON: %v", err)
		return err
	}

	userResult, err := b.userService.GetByTotalPointsAndContestID(ctx, args.Rating, data.ContestID)
	if err != nil {
		b.log.Error("userService.GetByTotalPointsAndContestID: %v", err)
		return err
	}

	rand.NewSource(time.Now().UnixNano())
	rand.Shuffle(len(userResult), func(i, j int) {
		userResult[i], userResult[j] = userResult[j], userResult[i]
	})

	// Выбираем первые n элементов после перемешивания
	selectedUsers := userResult[:args.UserNumber]

	usersByte, err := json.MarshalIndent(selectedUsers, "", "\t")
	if err != nil {
		b.log.Error("%v", err)
		return err
	}

	if err := b.tgMsg.SendNewMessage(update.FromChat().ID, nil, string(usersByte)); err != nil {
		b.log.Error("tgMsg.SendNewMessage in pickUsers:%v", err)
		return err
	}

	return nil
}

func (b *Bot) sendMessageToUser(ctx context.Context, update *tgbotapi.Update, data *store.ContestStore) error {
	args, err := ParseJSON[entity.ArgsUser](update.Message.Text)
	if err != nil {
		b.log.Error("ParseJSON: %v", err)
		return err
	}

	if err := b.tgMsg.SendNewMessage(args.UserID, nil, args.Message); err != nil {
		b.log.Error("tgMsg.SendNewMessage in pickUsers:%v", err)
		return err
	}

	if err := b.tgMsg.SendNewMessage(update.FromChat().ID, nil, "Сообщение пользователю отправлено"); err != nil {
		b.log.Error("tgMsg.SendNewMessage in pickUsers:%v", err)
		return err
	}

	return nil
}

func (b *Bot) updateRating(ctx context.Context, update *tgbotapi.Update, data *store.ContestStore) error {
	args, err := ParseJSON[entity.ArgsRating](update.Message.Text)
	if err != nil {
		b.log.Error("ParseJSON: %v", err)
		return err
	}

	tx, err := b.pg.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	if err := b.userService.UpdateTotalPointsByUserIDAndContestID(ctx, tx, args.UserID, data.ContestID, args.Rating); err != nil {
		b.log.Error("userService.UpdateTotalPointsByUserIDAndContestID: %v", err)
		return err
	}

	if err := b.tgMsg.SendNewMessage(update.FromChat().ID, nil, "Рейтинг изменен успешно"); err != nil {
		b.log.Error("tgMsg.SendNewMessage in pickUsers:%v", err)
		return err
	}

	return nil
}

func (b *Bot) additionQuestion(ctx context.Context, update *tgbotapi.Update, data *store.QuestionStore) error {
	args, err := ParseJSON[entity.ArgsTop10](update.Message.Text)
	if err != nil {
		b.log.Error("ParseJSON: %v", err)
		return err
	}
	args.ContestID = data.ContestID
	args.AdminID = data.UserID

	tx, err := b.pg.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	questionID, markupQuestion, err := b.answersService.CreateAdditionalQuestionWithAnswer(ctx, tx, args)
	if err != nil {
		b.log.Error("answersService.CreateAdditionalQuestionWithAnswer in additionQuestion: %v", err)
		return err
	}

	go func(ma *tgbotapi.InlineKeyboardMarkup, arg entity.ArgsTop10, questionID int) {
		var totalSend int

		for _, user := range arg.UsersID {

			if err := b.tgMsg.SendNewMessage(user, ma, arg.Question); err != nil {
				b.log.Error("tgMsg.SendNewMessage in user question send: %v", err)

				if strings.Contains(err.Error(), "Forbidden: bot was blocked by the user") ||
					strings.Contains(err.Error(), "Bad Request: chat not found") {

					if err := b.userService.UpdateBlockedBotStatus(context.Background(), user, true); err != nil {
						b.log.Error("userService.UpdateBlockedBotStatus: %v", err)
					}

				} else {
					b.log.Error("error on sending: %v", err)
				}
			}
			totalSend++
		}

		if err := b.questionsService.UpdateIsSendByQuestionID(context.Background(), true, questionID); err != nil {
			b.log.Error("questionsService.UpdateIsSendByQuestionID: %v", err)
		}

		if err := b.tgMsg.SendNewMessage(
			update.FromChat().ID,
			nil,
			fmt.Sprintf("Дополнительная рассылка завершена. Отправлено пользователям: %d", totalSend),
		); err != nil {
			b.log.Error("questionsService.UpdateIsSendByQuestionID: %v", err)
			return
		}
	}(markupQuestion, args, questionID)

	return nil
}

func (b *Bot) userMailing(ctx context.Context, update *tgbotapi.Update) error {
	args, err := ParseJSON[entity.ArgsMailing](update.Message.Text)
	if err != nil {
		b.log.Error("ParseJSON: %v", err)
		return err
	}

	user, err := b.userService.GetAllUsers(ctx)
	if err != nil {
		b.log.Error("userService.GetAllUsers: %v", err)
	}

	go func(user []entity.User, arg entity.ArgsMailing) {
		var totalSend int

		for _, user := range user {

			if user.BlockedBot == false {
				if err := b.tgMsg.SendNewMessage(user.ID, nil, arg.Message); err != nil {
					b.log.Error("tgMsg.SendNewMessage in user question send: %v", err)

					if strings.Contains(err.Error(), "Forbidden: bot was blocked by the user") ||
						strings.Contains(err.Error(), "Bad Request: chat not found") {

						if err := b.userService.UpdateBlockedBotStatus(context.Background(), user.ID, true); err != nil {
							b.log.Error("userService.UpdateBlockedBotStatus: %v", err)
						}

					} else {
						b.log.Error("error on sending: %v", err)
					}
				}
				totalSend++
			}
		}

		if err := b.tgMsg.SendNewMessage(
			update.FromChat().ID,
			nil,
			fmt.Sprintf("Рассылка завершена. Отправлено пользователям: %d", totalSend),
		); err != nil {
			b.log.Error("questionsService.UpdateIsSendByQuestionID: %v", err)
			return
		}
	}(user, args)

	return nil
}
