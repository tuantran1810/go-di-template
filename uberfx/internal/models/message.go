package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	Key   string
	Value string
}
