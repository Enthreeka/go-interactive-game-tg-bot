package tgbot

import (
	"context"
	"encoding/json"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/entity"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/store"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	"time"
)

type argsAddButton struct {
	NameAnswer   string `json:"ответ"`
	CostOfAnswer int    `json:"цена_ответа"`
}

type argsCreate struct {
	Name     string `json:"название_конкурса"`
	Deadline string `json:"дедлайн"`
}

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
		}

		if _, err := b.tgMsg.SendEditMessage(userID, update.Message.MessageID, nil, "Выполнено успешно"); err != nil {
			return true, err
		}

		return true, nil
	case *store.QuestionStore:
		b.log.Info("process store.QuestionStore")

		if data.TypeCommandQuestion == store.QuestionCreate {
			if err := b.questionsService.CreateQuestion(ctx, &entity.Question{
				ContestID:     data.ContestID,
				CreatedByUser: data.UserID,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
				QuestionName:  update.Message.Text,
			}); err != nil {
				b.log.Error("questionsService.CreateQuestion: %v", err)
				return true, err
			}
		}

		if data.TypeCommandQuestion == store.QuestionUpdate {
			if err := b.questionsService.UpdateQuestionName(ctx, data.QuestionID, update.Message.Text); err != nil {
				b.log.Error("questionsService.UpdateQuestionName: %v", err)
				return true, err
			}
		}

		if data.TypeCommandQuestion == store.QuestionAddButtonAnswer {
			if err := b.contentStoreAddAnswer(ctx, update, data.QuestionID); err != nil {
				b.log.Error("contentStoreAddAnswer: %v", err)
				return true, err
			}
		}

		if data.TypeCommandQuestion == store.QuestionAddDeadline {
			parsedTime, err := time.Parse(time.DateTime, update.Message.Text)
			if err != nil {
				return true, err
			}

			if err := b.questionsService.UpdateDeadlineByQuestionID(ctx, data.QuestionID, parsedTime); err != nil {
				b.log.Error("questionsService.UpdateDeadlineByQuestionID: %v", err)
				return true, err
			}
		}

		if err := b.updateChatInSuccessfullyCase(userID, data.MsgID, update.Message.MessageID); err != nil {
			return true, err
		}

		return true, nil
	default:
		b.log.Error("undefined type data in storeProcessing")
		return true, nil
	}
}

func (b *Bot) contentStoreAddAnswer(ctx context.Context, update *tgbotapi.Update, questionID int) error {
	args, err := ParseJSON[argsAddButton](update.Message.Text)
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
	args, err := ParseJSON[argsCreate](update.Message.Text)
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

func (b *Bot) updateChatInSuccessfullyCase(userID int64, telegramMessageID, userMessageID int) error {
	if resp, err := b.bot.Request(tgbotapi.NewDeleteMessage(userID,
		userMessageID)); nil != err || !resp.Ok {
		b.log.Error("failed to delete message id %d (%s): %v", userMessageID, string(resp.Result), err)
		return err
	}

	if _, err := b.tgMsg.SendEditMessage(userID, telegramMessageID, nil, "Выполнено успешно"); err != nil {
		return err
	}

	return nil
}
