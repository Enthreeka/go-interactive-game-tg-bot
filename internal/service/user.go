package service

import (
	"context"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/entity"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/repo/postgres"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/logger"
	"github.com/jackc/pgx/v5"
)

type UserService interface {
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
	CreateUserIfNotExist(ctx context.Context, user *entity.User) error
	GetAllUsers(ctx context.Context) ([]entity.User, error)
	UpdateBlockedBotStatus(ctx context.Context, userID int64, status bool) error

	IsExistUserResultByUserID(ctx context.Context, userID int64, contestID int) (bool, error)
	CreateUserResult(ctx context.Context, tx pgx.Tx, result *entity.UserResult) error
	GetUserResultsByContest(ctx context.Context, userID int64, contestID int) (*entity.UserResult, error)
	UpdateTotalPointsByUserIDAndContestID(ctx context.Context, tx pgx.Tx, userID int64, contestID int, totalPoint int) error
	UpdateTotalPointsByContestID(ctx context.Context, contestID int, totalPoint int) error
	GetByTotalPointsAndContestID(ctx context.Context, totalPoint, contestID int) ([]entity.UserResult, error)
	UpdateRoleByUsername(ctx context.Context, role string, username string) error

	GetTop10UserByContest(ctx context.Context, contestID int) ([]entity.UserResult, error)
	GetAllAdmin(ctx context.Context) ([]entity.User, error)
}

type userService struct {
	userRepo       postgres.UserRepo
	userResultRepo postgres.UserResultRepo
	log            *logger.Logger
}

func (u *userService) GetAllAdmin(ctx context.Context) ([]entity.User, error) {
	return u.userRepo.GetAllAdmin(ctx)
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

func (u *userService) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	return u.userRepo.GetAllUsers(ctx)
}

func (u *userService) UpdateBlockedBotStatus(ctx context.Context, userID int64, status bool) error {
	return u.userRepo.UpdateBlockedBotStatus(ctx, userID, status)
}

func (u *userService) IsExistUserResultByUserID(ctx context.Context, userID int64, contestID int) (bool, error) {
	return u.userResultRepo.IsExistUserResultByUserID(ctx, userID, contestID)
}

func (u *userService) CreateUserResult(ctx context.Context, tx pgx.Tx, result *entity.UserResult) error {
	return u.userResultRepo.CreateUserResult(ctx, tx, result)
}

func (u *userService) GetUserResultsByContest(ctx context.Context, userID int64, contestID int) (*entity.UserResult, error) {
	return u.userResultRepo.GetUserResultsByContest(ctx, userID, contestID)
}
func (u *userService) UpdateTotalPointsByUserIDAndContestID(ctx context.Context, tx pgx.Tx, userID int64, contestID int, totalPoint int) error {
	return u.userResultRepo.UpdateTotalPointsByUserIDAndContestID(ctx, tx, userID, contestID, totalPoint)
}

func (u *userService) UpdateTotalPointsByContestID(ctx context.Context, contestID int, totalPoint int) error {
	u.log.Info("UpdateTotalPointsByContestID: contestID = %d, totalPoint = %d", contestID, totalPoint)
	return u.userResultRepo.UpdateTotalPointsByContestID(ctx, contestID, totalPoint)
}

func (u *userService) GetByTotalPointsAndContestID(ctx context.Context, totalPoint, contestID int) ([]entity.UserResult, error) {
	return u.userResultRepo.GetByTotalPointsAndContestID(ctx, totalPoint, contestID)
}

func (u *userService) UpdateRoleByUsername(ctx context.Context, role string, username string) error {
	return u.userRepo.UpdateRoleByUsername(ctx, role, username)
}

func (u *userService) GetTop10UserByContest(ctx context.Context, contestID int) ([]entity.UserResult, error) {
	return u.userResultRepo.GetTop10UserByContest(ctx, contestID)
}
