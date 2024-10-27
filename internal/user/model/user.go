package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Name         string
	Email        string `gorm:"unique;not null"`
	Password     string `gorm:"-:migration"` // Don't create a `password` column in the database
	PasswordHash string `gorm:"not null"`
}
