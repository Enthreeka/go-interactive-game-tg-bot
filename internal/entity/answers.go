package entity

import (
	"strconv"
	"strings"
	"time"
)

type Answer struct {
	ID             int    `json:"id"`
	Answer         string `json:"answer"`
	CostOfResponse int    `json:"cost_of_response"`

	// -- helps field, no in DB
	QuestionID int        `json:"question_id,questions_id"`
	ContestID  int        `json:"contest_id,contest"`
	Deadline   *time.Time `json:"deadline"`
}

func GetAnswerID(data string) int {
	parts := strings.Split(data, "_")
	if len(parts) > 3 {
		return 0
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0
	}

	return id
}
