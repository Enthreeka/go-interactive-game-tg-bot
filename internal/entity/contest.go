package entity

import "time"

type Contest struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	FileID   string    `json:"file_id"`
	Deadline time.Time `json:"deadline,omitempty"`
}
