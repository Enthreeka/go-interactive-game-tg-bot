package entity

import (
	"strconv"
	"strings"
	"time"
)

type Contest struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	FileID      string    `json:"file_id"`
	Deadline    time.Time `json:"deadline,omitempty"`
	IsCompleted bool      `json:"is_completed"`
}

func GetContestID(data string) int {
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
