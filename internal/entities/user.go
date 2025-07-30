package entities

import (
	"time"
)

type User struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	Password  string
	Uuid      string
	Name      string
	Email     *string
}
