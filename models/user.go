package models

import "time"

type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Password  []byte    `json:"password,omitempty" db:"password"`
	Nickname  string    `json:"nickname" db:"nickname"`
	Avatar    string    `json:"avatar,omitempty" db:"avatar"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
