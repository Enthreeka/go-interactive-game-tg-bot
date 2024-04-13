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
)

type AnswersService interface {
	CreateAnswer(ctx context.Context, tx pgx.Tx, answer *entity.Answer) error
	GetAnswerByID(ctx context.Context, id int, method string) ([]entity.Answer, *tgbotapi.InlineKeyboardMarkup, error)
	DeleteAnswer(ctx context.Context, id int) error
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

func (a *answersService) GetAnswerByID(ctx context.Context, id int, method string) ([]entity.Answer, *tgbotapi.InlineKeyboardMarkup, error) {
	answers, err := a.answerRepo.GetAnswerByID(ctx, id)
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
