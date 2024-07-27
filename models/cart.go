package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	ID         uint      `json:"cart-id" gorm:"primaryKey;autoIncrement"`
	Total      float64   `json:"total" gorm:"type:float;not null"`
	IsPayed    bool      `json:"isPayed" gorm:"type:boolean;not null"`
	CustomerID uint      `json:"customer-id" gorm:"not null"`
	Customer   Customer  `gorm:"foreignKey:CustomerID"`
	Products   []Product `gorm:"many2many:cart_products"`
}

type CartProduct struct {
	CartID    uint `gorm:"primaryKey"`
	ProductID uint `gorm:"primaryKey"`
	Quantity  int  `json:"quantity"`
}
