package entity

type HistoryPoint struct {
	UserID       int64 `json:"user_id"`
	QuestionID   int   `json:"questions_id"`
	AwardedPoint int   `json:"awarded_point"`
}
