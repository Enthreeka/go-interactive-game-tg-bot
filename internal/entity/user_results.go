package entity

type UserResult struct {
	ID          int   `json:"id"`
	UserID      int64 `json:"user_id"`
	ContestID   int   `json:"contest_id,omitempty"`
	TotalPoints int   `json:"total_points"`

	User User `json:"user"`
}
