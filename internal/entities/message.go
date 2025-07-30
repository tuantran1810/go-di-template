package entities

import "time"

type Message struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Key       string
	Value     string
}
