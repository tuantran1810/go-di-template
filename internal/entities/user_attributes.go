package entities

import "time"

type KeyValuePair struct {
	Key   string
	Value string
}

type UserAttribute struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uint
	Key       string
	Value     string
}
