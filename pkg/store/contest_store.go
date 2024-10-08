package store

type TypeCommandContest string

var (
	ContestCreate TypeCommandContest = "create"
	ContestDelete TypeCommandContest = "delete"
	ContestPick   TypeCommandContest = "pick"
	ContestUser   TypeCommandContest = "user"
	ContestRating TypeCommandContest = "rating"

	MessageUpdate TypeCommandContest = "update_msg"

	CreateUserMailing TypeCommandContest = "mailing"
)

type ContestStore struct {
	MsgID  int
	UserID int64

	ContestID int

	TypeCommandContest TypeCommandContest
}
