package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Sessions []Session
}

type Session struct {
	gorm.Model
	Token     string    `gorm:"unique;not null;size:255"`
	ExpiresAt time.Time `gorm:"not null"`
	UserId    uint      `gorm:"not null;index"`
	User      User
}
