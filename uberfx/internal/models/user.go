package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex,size:32"`
	Password string
	Uuid     string
	Name     string
	Email    sql.NullString
}
