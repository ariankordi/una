package models

import (
	"time"

	"github.com/gobuffalo/pop/nulls"
)

type Lobby struct {
	ID                 int            `json:"id" db:"id"`
	Name               string         `json:"name" db:"name"`
	CreatorID          nulls.Int      `json:"creator_id,omitempty" db:"creator_id"`
	Creator            *User          `json:"creator,omitempty" db:"-"`
	AnonymousCreatorID nulls.Int      `json:"anonymous_creator_id,omitempty" db:"anonymous_creator_id"`
	AnonymousCreator   *AnonymousUser `json:"anonymous_user,omitempty" db:"-"`
	CreatedAt          time.Time      `json:"created_at" db:"created_at"`
}
