package postgres

import (
	"context"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/postgres"
)

type CommunicationRepo interface {
	Insert(ctx context.Context, text string) error
	Get(ctx context.Context) (string, error)
	Update(ctx context.Context, text string) error
	IsExist(ctx context.Context) (bool, error)
}

type communicationRepo struct {
	*postgres.Postgres
}

func NewCommunicationRepo(pg *postgres.Postgres) CommunicationRepo {
	return &communicationRepo{
		pg,
	}
}

func (c *communicationRepo) Insert(ctx context.Context, text string) error {
	query := `INSERT INTO communication (message) VALUES ($1)`
	_, err := c.Pool.Exec(ctx, query, text)
	return err
}

func (c *communicationRepo) Get(ctx context.Context) (string, error) {
	var text string
	query := `SELECT message FROM communication`
	err := c.Pool.QueryRow(ctx, query).Scan(&text)
	return text, err
}

func (c *communicationRepo) Update(ctx context.Context, text string) error {
	query := `UPDATE communication SET message = $1`
	_, err := c.Pool.Exec(ctx, query, text)
	return err
}

func (c *communicationRepo) IsExist(ctx context.Context) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT FROM communication)`
	err := c.Pool.QueryRow(ctx, query).Scan(&exists)
	return exists, err
}
