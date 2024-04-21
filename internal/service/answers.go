package service

import (
	"context"
	"fmt"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/entity"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/repo/postgres"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/logger"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/tg/button"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	"time"
)

type AnswersService interface {
	CreateAnswer(ctx context.Context, tx pgx.Tx, answer *entity.Answer) error
	GetAnswersByID(ctx context.Context, questionID int, method string) ([]entity.Answer, *tgbotapi.InlineKeyboardMarkup, error)
	DeleteAnswer(ctx context.Context, id int) error
	GetContestIDByQuestionID(ctx context.Context, questionID int) (int, error)
	GetAnswerByID(ctx context.Context, id int) (*entity.Answer, error)
	AddHistoryPoints(ctx context.Context, tx pgx.Tx, userID int64, questionID, awardedPoints int) error
	CreateAdditionalQuestionWithAnswer(ctx context.Context, tx pgx.Tx, args entity.ArgsTop10) (int, *tgbotapi.InlineKeyboardMarkup, error)

	Declension(number int) string
}

type answersService struct {
	answerRepo         postgres.AnswerRepo
	questionRepo       postgres.QuestionRepo
	questionAnswerRepo postgres.QuestionAnswerRepo
	historyPointsRepo  postgres.HistoryPointsRepo
	log                *logger.Logger
}

func NewAnswersService(answerRepo postgres.AnswerRepo,
	questionRepo postgres.QuestionRepo,
	questionAnswerRepo postgres.QuestionAnswerRepo,
	historyPointsRepo postgres.HistoryPointsRepo,
	log *logger.Logger) AnswersService {
	return &answersService{
		answerRepo:         answerRepo,
		questionRepo:       questionRepo,
		questionAnswerRepo: questionAnswerRepo,
		historyPointsRepo:  historyPointsRepo,
		log:                log,
	}
}

func (a *answersService) Declension(number int) string {
	switch {
	case number%10 == 1 && number%100 != 11:
		return "балл"
	case number%10 >= 2 && number%10 <= 4 && (number%100 < 10 || number%100 >= 20):
		return "балла"
	default:
		return "баллов"
	}
}

func (a *answersService) CreateAdditionalQuestionWithAnswer(ctx context.Context, tx pgx.Tx, args entity.ArgsTop10) (int, *tgbotapi.InlineKeyboardMarkup, error) {
	a.log.Info("Create new additional data: %v", args)

	questionID, err := a.questionRepo.CreateQuestion(ctx, tx, &entity.Question{
		ContestID:     args.ContestID,
		CreatedByUser: args.AdminID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		QuestionName:  args.Question,
		FileID:        nil,
	})
	if err != nil {
		a.log.Error("CreateAdditionalQuestionWithAnswer: questionRepo.CreateQuestion: %v", err)
		return 0, nil, err
	}
	var ans []entity.Answer
	for _, value := range args.Answers {
		var a entity.Answer
		a.Answer = value.Answer
		a.CostOfResponse = value.Cost
		ans = append(ans, a)
	}

	createdID, err := a.answerRepo.CreateAnswers(ctx, tx, ans)
	if err != nil {
		a.log.Error("CreateAdditionalQuestionWithAnswer: answerRepo.CreateAnswers: %v", err)
		return 0, nil, err
	}

	for _, answerID := range createdID {
		if err := a.questionAnswerRepo.LinkQuestionToAnswer(ctx, tx, questionID, answerID, args.ContestID); err != nil {
			a.log.Error("questionAnswerRepo.LinkQuestionToAnswer: %v", err)
			return 0, nil, err
		}
	}

	_, markup, err := a.GetAnswersByID(ctx, questionID, "get")
	if err != nil {
		a.log.Error("GetAnswersByID in CreateAdditionalQuestionWithAnswer: %v", err)
		return 0, nil, err
	}

	return questionID, markup, nil
}

func (a *answersService) CreateAnswer(ctx context.Context, tx pgx.Tx, answer *entity.Answer) error {
	a.log.Info("Create new answer: %v", answer)

	answerID, err := a.answerRepo.CreateAnswer(ctx, tx, answer)
	if err != nil {
		a.log.Error("answerRepo.CreateAnswer: %v", err)
		return err
	}

	contestID, err := a.questionRepo.GetContestIDByQuestionID(ctx, answer.QuestionID)
	if err != nil {
		a.log.Error("questionRepo.GetContestIDByQuestionID: %v", err)
		return err
	}

	if err := a.questionAnswerRepo.LinkQuestionToAnswer(ctx, tx, answer.QuestionID, answerID, contestID); err != nil {
		a.log.Error("questionAnswerRepo.LinkQuestionToAnswer: %v", err)
		return err
	}

	return nil
}

func (a *answersService) GetAnswersByID(ctx context.Context, questionID int, method string) ([]entity.Answer, *tgbotapi.InlineKeyboardMarkup, error) {
	answers, err := a.questionAnswerRepo.GetAnswersByQuestion(ctx, questionID)
	if err != nil {
		a.log.Error("answerRepo.GetAnswerByID: %v", err)
		return nil, nil, err
	}

	markup, err := a.createAnswerMarkup(answers, method)
	if err != nil {
		a.log.Error("createAnswerMarkup: %v", err)
		return nil, nil, err
	}

	return answers, markup, nil
}

func (q *answersService) createAnswerMarkup(answers []entity.Answer, method string) (*tgbotapi.InlineKeyboardMarkup, error) {
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	buttonsPerRow := 1

	for i, el := range answers {
		btn := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s", el.Answer),
			fmt.Sprintf("answer_%s_%d", method, el.ID))

		row = append(row, btn)

		if (i+1)%buttonsPerRow == 0 || i == len(answers)-1 {
			rows = append(rows, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}

	rows = append(rows, []tgbotapi.InlineKeyboardButton{button.MainMenuButton})
	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &markup, nil
}

func (a *answersService) DeleteAnswer(ctx context.Context, id int) error {
	return a.answerRepo.DeleteAnswer(ctx, id)
}

func (a *answersService) GetAnswerByID(ctx context.Context, id int) (*entity.Answer, error) {
	return a.answerRepo.GetAnswerByID(ctx, id)
}

func (a *answersService) GetContestIDByQuestionID(ctx context.Context, questionID int) (int, error) {
	return a.questionRepo.GetContestIDByQuestionID(ctx, questionID)
}

func (a *answersService) AddHistoryPoints(ctx context.Context, tx pgx.Tx, userID int64, questionID, awardedPoints int) error {
	a.log.Info("Insert new result: userID - %d, questionID - %d, awardedPoints - %d", userID, questionID, awardedPoints)
	return a.historyPointsRepo.AddHistoryPoints(ctx, tx, userID, questionID, awardedPoints)
}
