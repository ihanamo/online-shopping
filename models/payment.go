package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	ID         uint      `json:"payment-id" gorm:"primaryKey;autoIncrement"`
	PaidAt     time.Time `json:"paid_at" gorm:"not null"`
	Total      float64   `json:"total" gorm:"type:float;not null"`
	CustomerID uint      `json:"customer_id" gorm:"not null"`
	Customer   Customer  `gorm:"foreignKey:CustomerID"`
}
