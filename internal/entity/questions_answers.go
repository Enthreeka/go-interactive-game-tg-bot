package entity

type QuestionsAnswers struct {
	QuestionID int `json:"questions_id"`
	AnswerID   int `json:"answers_id"`
	ContestID  int `json:"contest_id"`
}
