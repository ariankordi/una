package models

import "time"

type AnonymousUser struct {
	ID        int       `json:"id" db:"id"`
	NameID    int       `json:"name_id" db:"name_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
