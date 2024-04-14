package postgres

import (
	"context"
	pgxError "github.com/Entreeka/go-interactive-game-tg-bot/internal/boterror/pgx_error"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/entity"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"time"
)

type QuestionRepo interface {
	GetQuestionsCreatedByUser(ctx context.Context, userID int64) ([]entity.Question, error)
	GetAllQuestions(ctx context.Context) ([]entity.Question, error)
	GetQuestionByID(ctx context.Context, id int) (*entity.Question, error)
	GetQuestionsByContestID(ctx context.Context, contestID int) ([]entity.Question, error)
	GetContestIDByQuestionID(ctx context.Context, questionID int) (int, error)

	UpdateQuestion(ctx context.Context, question *entity.Question) error
	UpdateQuestionName(ctx context.Context, questionID int, name string) error
	UpdateDeadlineByQuestionID(ctx context.Context, questionID int, deadline time.Time) error
	UpdateIsSendByQuestionID(ctx context.Context, isSend bool, questionID int) error

	DeleteQuestion(ctx context.Context, id int) error

	CreateQuestion(ctx context.Context, question *entity.Question) error
}

type questionRepo struct {
	*postgres.Postgres
}

func NewQuestionRepo(pg *postgres.Postgres) QuestionRepo {
	return &questionRepo{
		pg,
	}
}

func (q *questionRepo) collectRow(row pgx.Row) (*entity.Question, error) {
	var question entity.Question
	err := row.Scan(&question.ID,
		&question.ContestID,
		&question.CreatedByUser,
		&question.CreatedAt,
		&question.UpdatedAt,
		&question.QuestionName,
		&question.FileID,
		&question.Deadline,
		&question.IsSend,
	)
	if checkErr := pgxError.ErrorHandler(err); checkErr != nil {
		return nil, checkErr
	}

	return &question, err
}

func (q *questionRepo) collectRows(rows pgx.Rows) ([]entity.Question, error) {
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (entity.Question, error) {
		question, err := q.collectRow(row)
		return *question, err
	})
}

func (q *questionRepo) CreateQuestion(ctx context.Context, question *entity.Question) error {
	query := `INSERT INTO questions (created_by_user, created_at, updated_at, question_name, file_id,contest_id) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := q.Pool.Exec(ctx, query, question.CreatedByUser, question.CreatedAt, question.UpdatedAt, question.QuestionName, question.FileID, question.ContestID)
	return err
}

func (q *questionRepo) GetAllQuestions(ctx context.Context) ([]entity.Question, error) {
	query := `SELECT * FROM questions`

	rows, err := q.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	return q.collectRows(rows)
}

func (q *questionRepo) GetQuestionByID(ctx context.Context, id int) (*entity.Question, error) {
	query := `SELECT * FROM questions WHERE id = $1`

	row := q.Pool.QueryRow(ctx, query, id)
	return q.collectRow(row)
}

func (q *questionRepo) UpdateQuestion(ctx context.Context, question *entity.Question) error {
	query := `UPDATE questions SET updated_at = $1, question_name = $2, file_id = $3 WHERE id = $4`

	_, err := q.Pool.Exec(ctx, query, question.UpdatedAt, question.QuestionName, question.FileID, question.ID)
	return err
}

func (q *questionRepo) DeleteQuestion(ctx context.Context, id int) error {
	query := `DELETE FROM questions WHERE id = $1`

	_, err := q.Pool.Exec(ctx, query, id)
	return err
}

func (q *questionRepo) GetQuestionsCreatedByUser(ctx context.Context, userID int64) ([]entity.Question, error) {
	query := `SELECT * FROM questions WHERE created_by_user = $1`

	rows, err := q.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	return q.collectRows(rows)
}

func (q *questionRepo) GetQuestionsByContestID(ctx context.Context, contestID int) ([]entity.Question, error) {
	query := `SELECT  q.id,q.contest_id,q.created_by_user,q.created_at,q.updated_at,q.question_name,q.file_id,q.deadline,q.is_send from questions q
				join contest c on c.id = q.contest_id
				where c.id = $1`

	rows, err := q.Pool.Query(ctx, query, contestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []entity.Question
	for rows.Next() {
		var question entity.Question
		err := rows.Scan(&question.ID,
			&question.ContestID,
			&question.CreatedByUser,
			&question.CreatedAt,
			&question.UpdatedAt,
			&question.QuestionName,
			&question.FileID,
			&question.Deadline,
			&question.IsSend,
		)
		if err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}

	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}

func (q *questionRepo) UpdateQuestionName(ctx context.Context, questionID int, name string) error {
	query := `update questions set question_name = $1 where id = $2`

	_, err := q.Pool.Exec(ctx, query, name, questionID)
	return err
}

func (q *questionRepo) GetContestIDByQuestionID(ctx context.Context, questionID int) (int, error) {
	query := `select contest_id from questions where id = $1`
	var id int

	err := q.Pool.QueryRow(ctx, query, questionID).Scan(&id)
	return id, err
}

func (q *questionRepo) UpdateDeadlineByQuestionID(ctx context.Context, questionID int, deadline time.Time) error {
	query := `update questions set deadline = $1 where id = $2`

	_, err := q.Pool.Exec(ctx, query, deadline, questionID)
	return err
}

func (q *questionRepo) UpdateIsSendByQuestionID(ctx context.Context, isSend bool, questionID int) error {
	query := `update questions set is_send = $1 where id = $2`

	_, err := q.Pool.Exec(ctx, query, isSend, questionID)
	return err
}
