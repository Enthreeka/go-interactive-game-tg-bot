package entity

import "time"

type UserRole string

const (
	UserType       UserRole = "user"
	AdminType      UserRole = "admin"
	SuperAdminType UserRole = "superAdmin"
)

type User struct {
	ID          int64     `json:"id"`
	TGUsername  string    `json:"tg_username"`
	CreatedAt   time.Time `json:"created_at"`
	Phone       string    `json:"phone,omitempty"`
	ChannelFrom string    `json:"channel_from,omitempty"`
	UserRole    UserRole  `json:"user_role"`
}
