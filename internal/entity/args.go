package entity

type ArgsTop10 struct {
	Question string `json:"вопрос"`

	Answers []struct {
		Answer string `json:"ответ"`
		Cost   int    `json:"цена_ответа"`
	} `json:"варианты_ответы"`

	UsersID []int64 `json:"пользователи"`

	AdminID   int64
	ContestID int
}

type ArgsAddButton struct {
	NameAnswer   string `json:"ответ"`
	CostOfAnswer int    `json:"цена_ответа"`
}

type ArgsCreate struct {
	Name     string `json:"название_конкурса"`
	Deadline string `json:"дедлайн"`
}

type ArgsPick struct {
	Rating     int `json:"рейтинг"`
	UserNumber int `json:"количество_людей"`
}

type ArgsUser struct {
	Message string `json:"сообщение"`
	UserID  int64  `json:"id_пользователя"`
}

type ArgsRating struct {
	Rating int   `json:"рейтинг"`
	UserID int64 `json:"id_пользователя"`
}

type ArgsMailing struct {
	Message string `json:"сообщение"`
}
