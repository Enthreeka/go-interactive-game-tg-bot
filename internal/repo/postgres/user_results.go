package postgres

import (
	"context"
	"errors"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/entity"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type UserResultRepo interface {
	CreateUserResult(ctx context.Context, result *entity.UserResult) error
	GetUserResultsByContest(ctx context.Context, userID, contestID int) (*entity.UserResult, error)
	GetAllUserResultsByContest(ctx context.Context, contestID int) ([]entity.UserResult, error)
}

type userResultRepo struct {
	*postgres.Postgres
}

func NewUserResultRepo(pg *postgres.Postgres) UserResultRepo {
	return &userResultRepo{
		Postgres: pg,
	}
}

func (ur *userResultRepo) CreateUserResult(ctx context.Context, result *entity.UserResult) error {
	query := `INSERT INTO user_results (user_id, contest_id, total_points) VALUES ($1, $2, $3)`

	_, err := ur.Pool.Exec(ctx, query, result.UserID, result.ContestID, result.TotalPoints)
	return err
}

func (ur *userResultRepo) GetUserResultsByContest(ctx context.Context, userID, contestID int) (*entity.UserResult, error) {
	query := `SELECT id, total_points FROM user_results WHERE user_id = $1 AND contest_id = $2`

	row := ur.Pool.QueryRow(ctx, query, userID, contestID)
	var userResult entity.UserResult
	err := row.Scan(&userResult.ID, &userResult.TotalPoints)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, nil
	case err != nil:
		return nil, err
	}

	return &userResult, nil
}

func (ur *userResultRepo) GetAllUserResultsByContest(ctx context.Context, contestID int) ([]entity.UserResult, error) {
	query := `SELECT u.tg_username,user_results.user_id,user_results.id, user_results.total_points FROM user_results
                        join "user" u on u.id = user_results.user_id
                        WHERE user_results.contest_id = $1`

	rows, err := ur.Pool.Query(ctx, query, contestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []entity.UserResult
	for rows.Next() {
		var result entity.UserResult
		err := rows.Scan(&result.User.TGUsername, &result.UserID, &result.ID, &result.TotalPoints)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
