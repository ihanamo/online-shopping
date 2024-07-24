package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	FirstName string `json:"firstname" gorm:"type:varchar(100);not null"`
	LastName  string `json:"lastname" gorm:"type:varchar(100);not null"`
	Username  string `json:"username" gorm:"type:varchar(100);unique;not null"`
	Phone     string `json:"phone" gorm:"type:varchar(100);not null"`
	Password  string `json:"password" gorm:"type:varchar(255);not null"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	gorm.Model
	Token     string    `json:"token" gorm:"type:text;not null"`
	UserID    uint      `json:"userID" gorm:"not null"`
	ExpiresAt time.Time `json:"expiresAt" gorm:"not null"`
}
