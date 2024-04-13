package store

type TypeCommandContest string

var (
	ContestCreate TypeCommandContest = "create"
	ContestDelete TypeCommandContest = "delete"
)

type ContestStore struct {
	MsgID  int
	UserID int64

	TypeCommandContest TypeCommandContest
}
