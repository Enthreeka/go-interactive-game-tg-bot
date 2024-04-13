package service

import (
	"context"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/entity"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/repo/postgres"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/logger"
)

type UserService interface {
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
	CreateUserIfNotExist(ctx context.Context, user *entity.User) error
}

type userService struct {
	userRepo       postgres.UserRepo
	userResultRepo postgres.UserResultRepo
	log            *logger.Logger
}

func NewUserService(userRepo postgres.UserRepo, userResultRepo postgres.UserResultRepo, log *logger.Logger) UserService {
	return &userService{
		userRepo:       userRepo,
		userResultRepo: userResultRepo,
		log:            log,
	}
}

func (u *userService) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	return u.userRepo.GetUserByID(ctx, id)
}

func (u *userService) CreateUserIfNotExist(ctx context.Context, user *entity.User) error {
	isExist, err := u.userRepo.IsUserExistByUserID(ctx, user.ID)
	if err != nil {
		u.log.Error("userRepo.IsUserExistByUsernameTg: failed to check user: %v", err)
		return err
	}

	if !isExist {
		u.log.Info("Get user: %s", user.String())
		err := u.userRepo.CreateUser(ctx, user)
		if err != nil {
			u.log.Error("userRepo.CreateUser: failed to create user: %v", err)
			return err
		}
	}

	return nil
}
