package service

import (
	"context"
	"fmt"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/entity"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/repo/postgres"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/logger"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/tg/button"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

type QuestionsService interface {
	GetQuestionsByContestID(ctx context.Context, contestID int, method string) ([]entity.Question, *tgbotapi.InlineKeyboardMarkup, error)
	CreateQuestion(ctx context.Context, question *entity.Question) error
	GetQuestionByID(ctx context.Context, id int) (*entity.Question, error)
	UpdateQuestionName(ctx context.Context, questionID int, name string) error
	UpdateDeadlineByQuestionID(ctx context.Context, questionID int, deadline time.Time) error
}

type questionsService struct {
	questionRepo       postgres.QuestionRepo
	questionAnswerRepo postgres.QuestionAnswerRepo
	historyPointsRepo  postgres.HistoryPointsRepo
	log                *logger.Logger
}

func NewQuestionsService(questionRepo postgres.QuestionRepo, questionAnswerRepo postgres.QuestionAnswerRepo, historyPointsRepo postgres.HistoryPointsRepo, log *logger.Logger) QuestionsService {
	return &questionsService{
		questionRepo:       questionRepo,
		questionAnswerRepo: questionAnswerRepo,
		historyPointsRepo:  historyPointsRepo,
		log:                log,
	}
}

func (q *questionsService) CreateQuestion(ctx context.Context, question *entity.Question) error {
	q.log.Info("Create question: %#v", question)
	return q.questionRepo.CreateQuestion(ctx, question)
}

func (q *questionsService) GetQuestionsByContestID(ctx context.Context, contestID int, method string) ([]entity.Question, *tgbotapi.InlineKeyboardMarkup, error) {
	questions, err := q.questionRepo.GetQuestionsByContestID(ctx, contestID)
	if err != nil {
		q.log.Error("questionRepo.GetQuestionsByContestID: %v", err)
		return nil, nil, err
	}

	markup, err := q.createQuestionMarkup(questions, method, contestID)
	if err != nil {
		q.log.Error("createQuestionMarkup: %v", err)
		return nil, nil, err
	}

	return questions, markup, nil
}

func (q *questionsService) createQuestionMarkup(channel []entity.Question, method string, contestID int) (*tgbotapi.InlineKeyboardMarkup, error) {
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	buttonsPerRow := 1

	for i, el := range channel {
		btn := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s", el.QuestionName),
			fmt.Sprintf("question_%s_%d", method, el.ID))

		row = append(row, btn)

		if (i+1)%buttonsPerRow == 0 || i == len(channel)-1 {
			rows = append(rows, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}

	rows = append(rows, []tgbotapi.InlineKeyboardButton{button.BackToQuestionSetting(contestID)})
	rows = append(rows, []tgbotapi.InlineKeyboardButton{button.MainMenuButton})
	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &markup, nil
}

func (q *questionsService) GetQuestionByID(ctx context.Context, id int) (*entity.Question, error) {
	return q.questionRepo.GetQuestionByID(ctx, id)
}

func (q *questionsService) UpdateQuestionName(ctx context.Context, questionID int, name string) error {
	return q.questionRepo.UpdateQuestionName(ctx, questionID, name)
}

func (q *questionsService) UpdateDeadlineByQuestionID(ctx context.Context, questionID int, deadline time.Time) error {
	return q.questionRepo.UpdateDeadlineByQuestionID(ctx, questionID, deadline)
}
