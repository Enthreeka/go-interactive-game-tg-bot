package service

import (
	"context"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/repo/postgres"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/logger"
)

type CommunicationService interface {
	CreateMessage(ctx context.Context, message string) error
	GetMessage(ctx context.Context) (string, error)
}

type communicationService struct {
	log      *logger.Logger
	commRepo postgres.CommunicationRepo
}

func NewCommunicationService(commRepo postgres.CommunicationRepo, log *logger.Logger) CommunicationService {
	return &communicationService{
		commRepo: commRepo,
		log:      log,
	}
}

func (c *communicationService) CreateMessage(ctx context.Context, message string) error {
	isExistMsg, err := c.commRepo.IsExist(ctx)
	if err != nil {
		c.log.Error("Error while checking if communication is exist: %v", err)
		return err
	}

	if !isExistMsg {
		if err := c.commRepo.Insert(ctx, message); err != nil {
			c.log.Error("Error while inserting communication: %v", err)
			return err
		}
	}

	if err := c.commRepo.Update(ctx, message); err != nil {
		c.log.Error("Error while updating communication: %v", err)
		return err
	}

	return nil
}

func (c *communicationService) GetMessage(ctx context.Context) (string, error) {
	return c.commRepo.Get(ctx)
}
