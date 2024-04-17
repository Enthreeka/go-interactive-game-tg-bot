package postgres

import (
	"context"
	pgxError "github.com/Entreeka/go-interactive-game-tg-bot/internal/boterror/pgx_error"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/entity"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"time"
)

type ContestRepo interface {
	GetContestByID(ctx context.Context, id int) (*entity.Contest, error)
	GetAllContests(ctx context.Context) ([]entity.Contest, error)

	CreateContest(ctx context.Context, contest *entity.Contest) error

	UpdateContestDeadline(ctx context.Context, id int, deadline time.Time) error
	UpdateIsCompletedByContestID(ctx context.Context, isCompleted bool, contestID int) error

	DeleteContest(ctx context.Context, id int) error
}

type contestRepo struct {
	*postgres.Postgres
}

func NewContestRepo(pg *postgres.Postgres) ContestRepo {
	return &contestRepo{
		pg,
	}
}

func (c *contestRepo) collectRow(row pgx.Row) (*entity.Contest, error) {
	var contest entity.Contest
	err := row.Scan(&contest.ID, &contest.Name, &contest.FileID, &contest.Deadline, &contest.IsCompleted)
	if checkErr := pgxError.ErrorHandler(err); checkErr != nil {
		return nil, checkErr
	}

	return &contest, err
}

func (c *contestRepo) collectRows(rows pgx.Rows) ([]entity.Contest, error) {
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (entity.Contest, error) {
		contest, err := c.collectRow(row)
		return *contest, err
	})
}

func (c *contestRepo) GetContestByID(ctx context.Context, id int) (*entity.Contest, error) {
	query := `SELECT * FROM contest WHERE id = $1`

	row := c.Pool.QueryRow(ctx, query, id)
	return c.collectRow(row)
}

func (c *contestRepo) CreateContest(ctx context.Context, contest *entity.Contest) error {
	query := `INSERT INTO contest (name, file_id, deadline) VALUES ($1, $2, $3)`

	_, err := c.Pool.Exec(ctx, query, contest.Name, contest.FileID, contest.Deadline)
	return err
}

func (c *contestRepo) GetAllContests(ctx context.Context) ([]entity.Contest, error) {
	query := `SELECT * FROM contest`

	rows, err := c.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return c.collectRows(rows)
}

func (c *contestRepo) UpdateContestDeadline(ctx context.Context, id int, deadline time.Time) error {
	query := `UPDATE contest SET deadline = $1 WHERE id = $2`

	_, err := c.Pool.Exec(ctx, query, deadline, id)
	return err
}

func (c *contestRepo) DeleteContest(ctx context.Context, id int) error {
	query := `DELETE FROM contest WHERE id = $1`

	_, err := c.Pool.Exec(ctx, query, id)
	return err
}

func (c *contestRepo) UpdateIsCompletedByContestID(ctx context.Context, isCompleted bool, contestID int) error {
	query := `UPDATE contest SET is_completed = $1 where id = $2`

	_, err := c.Pool.Exec(ctx, query, isCompleted, contestID)
	return err
}
