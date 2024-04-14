package callback

import (
	"context"
	"errors"
	"fmt"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/boterror"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/entity"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/handler"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/handler/tgbot"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/service"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/logger"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/postgres"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/store"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/tg"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/tg/markup"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	"strings"
	"time"
)

type CallbackQuestion struct {
	questionsService service.QuestionsService
	answersService   service.AnswersService
	userService      service.UserService
	log              *logger.Logger
	store            *store.Store
	tgMsg            *tg.TelegramMsg
	pg               *postgres.Postgres
}

func NewCallbackQuestion(
	questionsService service.QuestionsService,
	answersService service.AnswersService,
	userService service.UserService,
	log *logger.Logger,
	store *store.Store,
	tgMsg *tg.TelegramMsg,
	pg *postgres.Postgres,
) *CallbackQuestion {
	return &CallbackQuestion{
		questionsService: questionsService,
		answersService:   answersService,
		userService:      userService,
		log:              log,
		store:            store,
		tgMsg:            tgMsg,
		pg:               pg,
	}
}

// CallbackQuestionSetting - question_setting_{contest_id}
func (c *CallbackQuestion) CallbackQuestionSetting() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		contestID := entity.GetContestID(update.CallbackData())

		markupQuestion := markup.QuestionSetting(contestID)
		if _, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&markupQuestion,
			"CallbackQuestionSetting"); err != nil {
			return err
		}

		return nil
	}
}

// CallbackGetAllQuestionByContestID - get_all_question_{contest_id}
func (c *CallbackQuestion) CallbackGetAllQuestionByContestID() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		contestID := entity.GetQuestionID(update.CallbackData())

		_, markupQuestion, err := c.questionsService.GetQuestionsByContestID(ctx, contestID, "get")
		if err != nil {
			c.log.Error("questionsService.GetQuestionsByContestID: failed to questions: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		if _, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			markupQuestion,
			"CallbackGetAllContest"); err != nil {
			return err
		}

		return nil
	}
}

// CallbackCreateQuestion - create_question_{contest_id}
func (c *CallbackQuestion) CallbackCreateQuestion() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		contestID := entity.GetContestID(update.CallbackData())

		msg, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&markup.CancelState,
			"Отправьте вопрос")
		if err != nil {
			return err
		}

		c.store.Delete(update.FromChat().ID)
		c.store.Set(&store.QuestionStore{
			MsgID:               msg,
			UserID:              update.FromChat().ID,
			ContestID:           contestID,
			TypeCommandQuestion: store.QuestionCreate,
		}, update.FromChat().ID)

		return nil
	}
}

// CallbackDeleteQuestion - delete_question
func (c *CallbackQuestion) CallbackDeleteQuestion() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		contestID := entity.GetContestID(update.CallbackData())

		_, markupQuestion, err := c.questionsService.GetQuestionsByContestID(ctx, contestID, "delete")
		if err != nil {
			c.log.Error("questionsService.GetQuestionsByContestID: failed to questions: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		if _, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			markupQuestion,
			"CallbackDeleteQuestion"); err != nil {
			return err
		}

		return nil
	}
}

// CallbackGetQuestionByID - question_get_{question_id}
func (c *CallbackQuestion) CallbackGetQuestionByID() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		questionID := entity.GetContestID(update.CallbackData())

		question, err := c.questionsService.GetQuestionByID(ctx, questionID)
		if err != nil {
			c.log.Error("questionsService.GetQuestionByID: failed to get question: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		questionMarkup := markup.QuestionByIDSetting(questionID, question.ContestID)
		if _, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&questionMarkup,
			fmt.Sprintf(question.QuestionName)); err != nil {
			return err
		}

		return nil
	}
}

// CallbackQuestionChangeName - question_change_name_{question_id}
func (c *CallbackQuestion) CallbackQuestionChangeName() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		questionID := entity.GetQuestionID(update.CallbackData())

		msg, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&markup.CancelState,
			"Отправьте измененный вопрос")
		if err != nil {
			return err
		}

		c.store.Delete(update.FromChat().ID)
		c.store.Set(&store.QuestionStore{
			MsgID:               msg,
			UserID:              update.FromChat().ID,
			QuestionID:          questionID,
			TypeCommandQuestion: store.QuestionUpdate,
		}, update.FromChat().ID)

		return nil
	}
}

// CallbackQuestionAddAnswer - question_add_answer_{question_id}
func (c *CallbackQuestion) CallbackQuestionAddAnswer() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		questionID := entity.GetQuestionID(update.CallbackData())

		msg, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&markup.CancelState,
			"{\n  \"ответ\": \"впишите сюда ответ\",\n  \"цена_ответа\": цену нужно указывать без скобок, целое число\n}")
		if err != nil {
			return err
		}

		c.store.Delete(update.FromChat().ID)
		c.store.Set(&store.QuestionStore{
			MsgID:               msg,
			UserID:              update.FromChat().ID,
			QuestionID:          questionID,
			TypeCommandQuestion: store.QuestionAddButtonAnswer,
		}, update.FromChat().ID)

		return nil
	}
}

// CallbackQuestionDeleteAnswer - question_delete_answer_{question_id}
func (c *CallbackQuestion) CallbackQuestionDeleteAnswer() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		questionID := entity.GetQuestionID(update.CallbackData())

		_, markupAnswer, err := c.answersService.GetAnswersByID(ctx, questionID, "delete")
		if err != nil {
			c.log.Error("answersService.GetAnswerByID: failed to get answer: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		if _, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			markupAnswer,
			"CallbackQuestionDeleteAnswer"); err != nil {
			return err
		}

		return nil
	}
}

// CallbackAnswerDelete - answer_delete_{answer_id}
func (c *CallbackQuestion) CallbackAnswerDelete() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		answerID := entity.GetAnswerID(update.CallbackData())

		if err := c.answersService.DeleteAnswer(ctx, answerID); err != nil {
			c.log.Error("answersService.DeleteAnswer: failed to delete answer: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		if _, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			nil,
			"Удалено успешно"); err != nil {
			return err
		}

		return nil
	}
}

// CallbackQuestionChangeDeadline - question_change_deadline_{question_id}
func (c *CallbackQuestion) CallbackQuestionChangeDeadline() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		questionID := entity.GetQuestionID(update.CallbackData())

		msg, err := c.tgMsg.SendEditMessage(update.FromChat().ID,
			update.CallbackQuery.Message.MessageID,
			&markup.CancelState,
			"Отправьте дедлайн")
		if err != nil {
			return err
		}

		c.store.Delete(update.FromChat().ID)
		c.store.Set(&store.QuestionStore{
			MsgID:               msg,
			UserID:              update.FromChat().ID,
			ContestID:           questionID,
			TypeCommandQuestion: store.QuestionAddDeadline,
		}, update.FromChat().ID)

		return nil
	}
}

// CallbackQuestionAdminView - question_admin_view_{question_id}
func (c *CallbackQuestion) CallbackQuestionAdminView() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		questionID := entity.GetQuestionID(update.CallbackData())

		question, err := c.questionsService.GetQuestionByID(ctx, questionID)
		if err != nil {
			c.log.Error("questionsService.GetQuestionByID: failed to get question: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		_, markupAnswer, err := c.questionsService.GetAnswersByQuestion(ctx, questionID, "get")
		if err != nil {
			if errors.Is(err, boterror.ErrEmptyAnswer) {
				handler.HandleError(bot, update, fmt.Sprintf("Необходимо создать ответы - [%s]", question.QuestionName))
				return nil
			}
			c.log.Error("questionsService.GetAnswersByQuestion: failed to get markup answers: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		if err := c.tgMsg.SendNewMessage(update.FromChat().ID, markupAnswer, question.QuestionName); err != nil {
			return err
		}

		return nil
	}
}

// CallbackQuestionSendUser - question_send_user_{question_id}
func (c *CallbackQuestion) CallbackQuestionSendUser() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		questionID := entity.GetQuestionID(update.CallbackData())

		question, err := c.questionsService.GetQuestionByID(ctx, questionID)
		if err != nil {
			c.log.Error("questionsService.GetQuestionByID: failed to get question: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		_, markupAnswer, err := c.questionsService.GetAnswersByQuestion(ctx, questionID, "get")
		if err != nil {
			if errors.Is(err, boterror.ErrEmptyAnswer) {
				handler.HandleError(bot, update, fmt.Sprintf("Необходимо создать ответы - [%s]", question.QuestionName))
				return nil
			}
			c.log.Error("questionsService.GetAnswersByQuestion: failed to get markup answers: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		if question == nil || question.QuestionName == "" {
			c.log.Error("question == nil || question.QuestionName == ''")
			handler.HandleError(bot, update, "Отсутствует вопрос")
			return nil
		}

		users, err := c.userService.GetAllUsers(ctx)
		if err != nil {
			c.log.Error("userService.GetAllUsers: failed to get users: %v", err)
			handler.HandleError(bot, update, boterror.ParseErrToText(err))
			return nil
		}

		go func(ma *tgbotapi.InlineKeyboardMarkup, u []entity.User, q *entity.Question, adminID int64) {
			var totalSend int

			for _, user := range u {

				if user.BlockedBot == false {
					if err := c.tgMsg.SendNewMessage(user.ID, markupAnswer, question.QuestionName); err != nil {
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

			if err := c.questionsService.UpdateIsSendByQuestionID(context.Background(), true, q.ID); err != nil {
				c.log.Error("questionsService.UpdateIsSendByQuestionID: %v", err)
			}

			if err := c.tgMsg.SendNewMessage(
				adminID,
				nil,
				fmt.Sprintf("Рассылка завершена. Отправлено пользователям: %d", totalSend),
			); err != nil {
				c.log.Error("questionsService.UpdateIsSendByQuestionID: %v", err)
				return
			}
		}(markupAnswer, users, question, update.FromChat().ID)

		return nil
	}
}

// CallbackAnswerGet - answer_get_{answer_id}
func (c *CallbackQuestion) CallbackAnswerGet() tgbot.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		var (
			awardedPoints int
			timeIsUp      string
			answerID      = entity.GetAnswerID(update.CallbackData())
			userID        = update.FromChat().ID
		)

		answer, err := c.answersService.GetAnswerByID(ctx, answerID)
		if err != nil {
			c.log.Error("answersService.GetAnswerByID: failed to get answer: %v", err)
			return nil
		}

		if answer.Deadline != nil {
			loc, _ := time.LoadLocation("Europe/Moscow")
			currentTime := time.Now().In(loc)

			if !currentTime.After(*answer.Deadline) {
				awardedPoints = answer.CostOfResponse
			} else {
				timeIsUp = "К сожалению, вы ответили после установленных сроков и не заработали баллов."
			}
		} else {
			awardedPoints = answer.CostOfResponse
		}

		tx, err := c.pg.Pool.BeginTx(ctx, pgx.TxOptions{})
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

		isExist, err := c.userService.IsExistUserResultByUserID(ctx, userID, answer.ContestID)
		if err != nil {
			c.log.Error("userService.IsExistUserResultByUserID: %v", err)
			return nil
		}
		if !isExist {
			if err := c.userService.CreateUserResult(ctx, tx, &entity.UserResult{
				UserID:      userID,
				ContestID:   answer.ContestID,
				TotalPoints: awardedPoints,
			}); err != nil {
				c.log.Error("userService.CreateUserResult: %v", err)
				return nil
			}
		} else {
			userResult, err := c.userService.GetUserResultsByContest(ctx, userID, answer.ContestID)
			if err != nil {
				c.log.Error("userService.GetUserResultsByContest: %v", err)
				return nil
			}

			if err := c.userService.UpdateTotalPointsByUserIDAndContestID(ctx, tx, userID, answer.ContestID, userResult.TotalPoints+awardedPoints); err != nil {
				c.log.Error("userService.UpdateTotalPointsByUserIDAndContestID: %v", err)
				return nil
			}
		}

		if err := c.answersService.AddHistoryPoints(ctx, tx, userID, answer.QuestionID, awardedPoints); err != nil {
			c.log.Error("answersService.AddHistoryPoints: failed insert in history_answer: %v", err)
			return nil
		}

		if _, err := c.tgMsg.SendEditMessage(
			userID,
			update.CallbackQuery.Message.MessageID,
			nil,
			fmt.Sprintf("Спасибо за ответ! Вы заработали: %d балла. %s", awardedPoints, timeIsUp),
		); err != nil {
			return nil
		}

		return nil
	}
}
