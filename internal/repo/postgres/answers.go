package postgres

import (
	"context"
	pgxError "github.com/Entreeka/go-interactive-game-tg-bot/internal/boterror/pgx_error"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/entity"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type AnswerRepo interface {
	GetAllAnswers(ctx context.Context) ([]entity.Answer, error)
	GetAnswerByID(ctx context.Context, id int) (*entity.Answer, error)

	CreateAnswer(ctx context.Context, tx pgx.Tx, answer *entity.Answer) (int, error)

	UpdateAnswer(ctx context.Context, answer *entity.Answer) error

	DeleteAnswer(ctx context.Context, id int) error
}

type answerRepo struct {
	*postgres.Postgres
}

func NewAnswerRepo(pg *postgres.Postgres) AnswerRepo {
	return &answerRepo{
		pg,
	}
}

func (a *answerRepo) collectRow(row pgx.Row) (*entity.Answer, error) {
	var answer entity.Answer
	err := row.Scan(&answer.ID, &answer.Answer, &answer.CostOfResponse, &answer.QuestionID, &answer.Deadline, &answer.ContestID)
	if checkErr := pgxError.ErrorHandler(err); checkErr != nil {
		return nil, checkErr
	}

	return &answer, err
}

func (a *answerRepo) collectRows(rows pgx.Rows) ([]entity.Answer, error) {
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (entity.Answer, error) {
		answer, err := a.collectRow(row)
		return *answer, err
	})
}

func (a *answerRepo) CreateAnswer(ctx context.Context, tx pgx.Tx, answer *entity.Answer) (int, error) {
	query := `INSERT INTO answers (answer, cost_of_response) VALUES ($1, $2) RETURNING id`
	var id int

	err := tx.QueryRow(ctx, query, answer.Answer, answer.CostOfResponse).Scan(&id)
	return id, err
}

func (a *answerRepo) GetAllAnswers(ctx context.Context) ([]entity.Answer, error) {
	query := `SELECT * FROM answers`

	rows, err := a.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return a.collectRows(rows)
}

func (a *answerRepo) GetAnswerByID(ctx context.Context, id int) (*entity.Answer, error) {
	query := `SELECT a.id, a.answer, a.cost_of_response, qa.questions_id, q.deadline, qa.contest FROM answers a
				join public.questions_answers qa on a.id = qa.answers_id
				join public.questions q on q.id = qa.questions_id
			WHERE a.id = $1
`

	row := a.Pool.QueryRow(ctx, query, id)
	return a.collectRow(row)
}

func (a *answerRepo) UpdateAnswer(ctx context.Context, answer *entity.Answer) error {
	query := `UPDATE answers SET answer = $1, cost_of_response = $2 WHERE id = $3`

	_, err := a.Pool.Exec(ctx, query, answer.Answer, answer.CostOfResponse, answer.ID)
	return err
}

func (a *answerRepo) DeleteAnswer(ctx context.Context, id int) error {
	query := `DELETE FROM answers WHERE id = $1`

	_, err := a.Pool.Exec(ctx, query, id)
	return err
}
