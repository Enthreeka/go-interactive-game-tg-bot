package entity

import "time"

type Question struct {
	ID           int       `json:"id"`
	CreatedBy    int64     `json:"created_by_user"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	QuestionName string    `json:"question_name"`
	FileID       string    `json:"file_id"`
}
