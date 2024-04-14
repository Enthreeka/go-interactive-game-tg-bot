package postgres

import (
	"context"
	"errors"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type HistoryPointsRepo interface {
	AddHistoryPoints(ctx context.Context, tx pgx.Tx, userID int64, questionID, awardedPoints int) error
	GetHistoryPoints(ctx context.Context, userID, questionID int) (int, error)
}

type historyPointsRepo struct {
	*postgres.Postgres
}

func NewHistoryPointsRepo(pg *postgres.Postgres) HistoryPointsRepo {
	return &historyPointsRepo{
		Postgres: pg,
	}
}

func (hp *historyPointsRepo) AddHistoryPoints(ctx context.Context, tx pgx.Tx, userID int64, questionID, awardedPoints int) error {
	query := `INSERT INTO history_points (user_id, questions_id, awarded_point) VALUES ($1, $2, $3)`

	_, err := tx.Exec(ctx, query, userID, questionID, awardedPoints)
	return err
}

func (hp *historyPointsRepo) GetHistoryPoints(ctx context.Context, userID, questionID int) (int, error) {
	query := `SELECT awarded_point FROM history_points WHERE user_id = $1 AND questions_id = $2`

	row := hp.Pool.QueryRow(ctx, query, userID, questionID)
	var awardedPoints int
	err := row.Scan(&awardedPoints)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return 0, nil
	case err != nil:
		return 0, err
	}

	return awardedPoints, nil
}
