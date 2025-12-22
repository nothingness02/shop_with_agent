package Order

import (
	"time"

	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusRefunded  OrderStatus = "refunded"
)

type Order struct {
	gorm.Model
	OrderID        string      `gorm:"uniqueIndex;size:64"`
	UserID         uint        `gorm:"index"`
	TotalAmount    float64     `gorm:"type:decimal(10,2)"`
	DiscountAmount float64     `gorm:"type:decimal(10,2)"`
	ShippingFee    float64     `gorm:"type:decimal(10,2)"`
	ActualAmount   float64     `gorm:"type:decimal(10,2)"`
	Status         OrderStatus `gorm:"size:32;index"`
	ExpiresAt      time.Time   `gorm:"type:date"`

	ShippingName    string `gorm:"size:100"`
	ShippingPhone   string `gorm:"size:20"`
	ShippingAddress string `gorm:"size:255"`
	ShippingZipCode string `gorm:"size:10"`
	Tsv             string `gorm:"type:tsvector;index:,type:gin;->"`

	OrderItems []OrderItem `gorm:"foreignKey:OrderID;references:ID"`
}

type OrderItem struct {
	gorm.Model
	OrderID     uint    `gorm:"index"`
	ProductID   uint    `gorm:"index"`
	ProductName string  `gorm:"size:100"`
	ProductImg  string  `gorm:"size:500"`
	Price       float64 `gorm:"type:decimal(10,2)"`
	Quantity    int
	Subtotal    float64 `gorm:"type:decimal(10,2)"`
}