package store

type TypeCommand string

var (
	UserAdminCreate TypeCommand = "create"
	UserAdminDelete TypeCommand = "delete"
)

type AdminStore struct {
	MsgID  int64
	UserID int64

	TypeCommand TypeCommand
}

func (a AdminStore) GetMsgID() int64 {
	return a.MsgID
}

func (a AdminStore) GetUserID() int64 {
	return a.UserID
}

func (a AdminStore) GetTypeCommand() TypeCommand {
	return a.TypeCommand
}
