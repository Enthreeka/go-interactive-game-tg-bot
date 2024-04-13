package store

type TypeCommandQuestion string

var (
	QuestionCreate          TypeCommandQuestion = "create"
	QuestionDelete          TypeCommandQuestion = "delete"
	QuestionUpdate          TypeCommandQuestion = "update"
	QuestionAddButtonAnswer TypeCommandQuestion = "add_button"
	QuestionAddDeadline     TypeCommandQuestion = "update_deadline"
)

type QuestionStore struct {
	MsgID      int
	UserID     int64
	ContestID  int
	QuestionID int

	TypeCommandQuestion TypeCommandQuestion
}
