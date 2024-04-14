package postgres

import (
	"context"
	"github.com/Entreeka/go-interactive-game-tg-bot/internal/entity"
	"github.com/Entreeka/go-interactive-game-tg-bot/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type QuestionAnswerRepo interface {
	LinkQuestionToAnswer(ctx context.Context, tx pgx.Tx, questionID, answerID, contestID int) error

	GetQuestionsByContest(ctx context.Context, contestID int) ([]entity.Question, error)
	GetAnswersByContest(ctx context.Context, contestID int) ([]entity.Answer, error)
	GetAnswersByQuestion(ctx context.Context, questionID int) ([]entity.Answer, error)
}

type questionAnswerRepo struct {
	*postgres.Postgres
}

func NewQuestionAnswerRepo(pg *postgres.Postgres) QuestionAnswerRepo {
	return &questionAnswerRepo{
		Postgres: pg,
	}
}

func (qa *questionAnswerRepo) LinkQuestionToAnswer(ctx context.Context, tx pgx.Tx, questionID, answerID, contestID int) error {
	query := `INSERT INTO questions_answers (questions_id, answers_id, contest) VALUES ($1, $2, $3)`

	_, err := tx.Exec(ctx, query, questionID, answerID, contestID)
	return err
}

func (qa *questionAnswerRepo) GetQuestionsByContest(ctx context.Context, contestID int) ([]entity.Question, error) {
	query := `SELECT q.id, q.created_by_user, q.created_at, q.updated_at, q.question_name, q.file_id
	          FROM questions q
	          INNER JOIN questions_answers qa ON q.id = qa.questions_id
	          WHERE qa.contest = $1`

	rows, err := qa.Pool.Query(ctx, query, contestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []entity.Question
	for rows.Next() {
		var question entity.Question
		err := rows.Scan(&question.ID, &question.CreatedByUser, &question.CreatedAt, &question.UpdatedAt, &question.QuestionName, &question.FileID)
		if err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}

func (qa *questionAnswerRepo) GetAnswersByContest(ctx context.Context, contestID int) ([]entity.Answer, error) {
	query := `SELECT a.id, a.answer, a.cost_of_response
	          FROM answers a
	          INNER JOIN questions_answers qa ON a.id = qa.answers_id
	          WHERE qa.contest = $1`

	rows, err := qa.Pool.Query(ctx, query, contestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []entity.Answer
	for rows.Next() {
		var answer entity.Answer
		err := rows.Scan(&answer.ID, &answer.Answer, &answer.CostOfResponse)
		if err != nil {
			return nil, err
		}
		answers = append(answers, answer)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return answers, nil
}

func (qa *questionAnswerRepo) GetAnswersByQuestion(ctx context.Context, questionID int) ([]entity.Answer, error) {
	query := `SELECT a.id, a.answer, a.cost_of_response
	          FROM answers a
	          INNER JOIN questions_answers qa ON a.id = qa.answers_id
	          WHERE qa.questions_id = $1`

	rows, err := qa.Pool.Query(ctx, query, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []entity.Answer
	for rows.Next() {
		var answer entity.Answer
		err := rows.Scan(&answer.ID, &answer.Answer, &answer.CostOfResponse)
		if err != nil {
			return nil, err
		}
		answers = append(answers, answer)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return answers, nil
}
