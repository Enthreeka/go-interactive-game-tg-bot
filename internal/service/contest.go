package service

import (
	"context"
	"fmt"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/entity"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/repo/postgres"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/logger"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/tg/button"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ContestService interface {
	GetAllContestsButtons(ctx context.Context, method string) ([]entity.Contest, *tgbotapi.InlineKeyboardMarkup, error)
	CreateContest(ctx context.Context, contest *entity.Contest) error
	GetContestByID(ctx context.Context, id int) (*entity.Contest, error)
	GetAllUserResultsByContest(ctx context.Context, contestID int) ([]entity.UserResult, error)
	DeleteContest(ctx context.Context, id int) error
}

type contestService struct {
	contestRepo    postgres.ContestRepo
	userResultRepo postgres.UserResultRepo
	log            *logger.Logger
}

func NewContestService(contestRepo postgres.ContestRepo, userResultRepo postgres.UserResultRepo, log *logger.Logger) ContestService {
	return &contestService{
		contestRepo:    contestRepo,
		userResultRepo: userResultRepo,
		log:            log,
	}
}

func (c *contestService) GetAllContestsButtons(ctx context.Context, method string) ([]entity.Contest, *tgbotapi.InlineKeyboardMarkup, error) {
	contest, err := c.contestRepo.GetAllContests(ctx)
	if err != nil {
		c.log.Error("GetAllContestsButtons: contestRepo.GetAllContests: %v", err)
		return nil, nil, err
	}

	markup, err := c.createContestMarkup(contest, method)
	if err != nil {
		c.log.Error("failed to create markup: %v", err)
		return nil, nil, err
	}

	return contest, markup, nil
}

func (c *contestService) CreateContest(ctx context.Context, contest *entity.Contest) error {
	return c.contestRepo.CreateContest(ctx, contest)
}

func (c *contestService) createContestMarkup(channel []entity.Contest, method string) (*tgbotapi.InlineKeyboardMarkup, error) {
	var rows [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	buttonsPerRow := 1

	for i, el := range channel {
		btn := tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s - %v", el.Name, el.Deadline),
			fmt.Sprintf("contest_%s_%d", method, el.ID))

		row = append(row, btn)

		if (i+1)%buttonsPerRow == 0 || i == len(channel)-1 {
			rows = append(rows, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}

	rows = append(rows, []tgbotapi.InlineKeyboardButton{button.BackToContestSetting})
	rows = append(rows, []tgbotapi.InlineKeyboardButton{button.MainMenuButton})
	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &markup, nil
}

func (c *contestService) GetContestByID(ctx context.Context, id int) (*entity.Contest, error) {
	return c.contestRepo.GetContestByID(ctx, id)
}

func (c *contestService) GetAllUserResultsByContest(ctx context.Context, contestID int) ([]entity.UserResult, error) {
	return c.userResultRepo.GetAllUserResultsByContest(ctx, contestID)
}

func (c *contestService) DeleteContest(ctx context.Context, id int) error {
	return c.contestRepo.DeleteContest(ctx, id)
}
