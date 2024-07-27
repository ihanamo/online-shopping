package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ID    uint    `json:"product-id" gorm:"primaryKey;autoIncrement"`
	Name  string  `json:"name" gorm:"type:varchar(100);not null"`
	Type  string  `json:"type" gorm:"type:varchar(100);not null"`
	Price float64 `json:"price" gorm:"type:float;not null"`
	Stock int     `json:"stock" gorm:"type:int;not null"`
	Carts []Cart  `gorm:"many2many:cart_products;"`
}
