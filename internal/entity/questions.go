package entity

import (
	"strconv"
	"strings"
	"time"
)

type Question struct {
	ID            int        `json:"id"`
	ContestID     int        `json:"contest_id"`
	CreatedByUser int64      `json:"created_by_user"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	QuestionName  string     `json:"question_name"`
	Deadline      *time.Time `json:"deadline"`
	FileID        *string    `json:"file_id"`
}

func GetQuestionID(data string) int {
	parts := strings.Split(data, "_")
	if len(parts) > 4 {
		return 0
	}

	id, err := strconv.Atoi(parts[3])
	if err != nil {
		return 0
	}

	return id
}
