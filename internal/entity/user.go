package entity

import (
	"fmt"
	"time"
)

type UserRole string

const (
	UserType       UserRole = "user"
	AdminType      UserRole = "admin"
	SuperAdminType UserRole = "superAdmin"
)

type User struct {
	ID          int64     `json:"id,omitempty"`
	TGUsername  string    `json:"tg_username"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	Phone       string    `json:"phone,omitempty"`
	ChannelFrom string    `json:"channel_from,omitempty"`
	UserRole    UserRole  `json:"user_role,omitempty"`
	BlockedBot  bool      `json:"blocked_bot,omitempty"`
}

func (u User) String() string {
	return fmt.Sprintf("(id: %d | tg_username: %s | channel_from: %v | created_at: %v | role: %s)",
		u.ID, u.TGUsername, u.ChannelFrom, u.CreatedAt, u.UserRole)
}
