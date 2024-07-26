package models

import (
	"time"

	"gorm.io/gorm"
)

type Manager struct {
	ID       uint   `json:"manager-id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name" gorm:"type:varchar(100);not null"`
	Username string `json:"username" gorm:"type:varchar(100);not null"`
	Password string `json:"password" gorm:"type:varchar(255);not null"`
}

type LoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JWTToken struct {
	gorm.Model
	Token     string    `json:"token" gorm:"type:text;not null"`
	UserID    uint      `json:"userID" gorm:"not null"`
	ExpiresAt time.Time `json:"expiresAt" gorm:"not null"`
}
