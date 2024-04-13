package callback

import (
	"context"
	"fmt"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/boterror"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/entity"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/handler"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/handler/tgbot"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/service"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/logger"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/store"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/tg"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/tg/markup"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackQuestion struct {
	questionsService service.QuestionsService
	answersService   service.AnswersService
	log              *logger.Logger
	store            *store.Store
	tgMsg            *tg.TelegramMsg
}

func NewCallbackQuestion(questionsService service.QuestionsService, answersService service.AnswersService, log *logger.Logger, store *store.Store, tgMsg *tg.TelegramMsg) *CallbackQuestion {
	return &CallbackQuestion{
		questionsService: questionsService,
		answersService:   answersService,
		log:              log,
		store:            store,
		tgMsg:            tgMsg,
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

		_, markupAnswer, err := c.answersService.GetAnswerByID(ctx, questionID, "delete")
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
