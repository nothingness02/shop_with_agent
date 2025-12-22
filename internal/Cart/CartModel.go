package cart

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	UserID      uint    `gorm:"index"`
	ProductID   uint    `gorm:"index"`
	ProductName string  `gorm:"size:100"`
	ProductImg  string  `gorm:"size:500"`
	Price       float64 `gorm:"type:decimal(10,2)"`
	Quantity    int
}