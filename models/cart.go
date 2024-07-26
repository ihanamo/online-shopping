package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	ID         uint     `json:"cart-id" gorm:"primaryKey;autoIncrement"`
	Total      float64  `json:"total" gorm:"type:float;not null"`
	IsPayed    bool     `json:"isPayed" gorm:"type:boolean;not null"`
	ProductID  uint     `json:"product-id"`
	Product    Product  `gorm:"foreignKey:ProductID"`
	CustomerID Customer `json:"customer-id" gorm:"not null"`
	Customer   Customer `gorm:"foreignKey:CustomerID"`
}
