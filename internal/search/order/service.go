package ordersearch

import (
	"fmt"

	"gorm.io/gorm"
)

// Service handles order full-text search.
type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// Result is a lightweight view returned to callers.
type Result struct {
	ID              uint    `json:"id"`
	OrderID         string  `json:"order_id"`
	UserID          uint    `json:"user_id"`
	Status          string  `json:"status"`
	ShippingName    string  `json:"shipping_name"`
	ShippingPhone   string  `json:"shipping_phone"`
	ShippingAddress string  `json:"shipping_address"`
	ShippingZipCode string  `json:"shipping_zip_code"`
	TotalAmount     float64 `json:"total_amount"`
}

// orderRecord mirrors the database schema for search queries.
type orderRecord struct {
	gorm.Model
	OrderID         string
	UserID          uint
	Status          string
	ShippingName    string
	ShippingPhone   string
	ShippingAddress string
	ShippingZipCode string
	TotalAmount     float64
	Tsv             string `gorm:"type:tsvector;index:,type:gin;->"`
}

func (orderRecord) TableName() string {
	return "orders"
}

// Search runs a full-text search against orders.
func (s *Service) Search(query, lang string) ([]Result, error) {
	var rows []orderRecord

	tsquery := gorm.Expr("plainto_tsquery(?, ?)", lang, query)
	tsvector := gorm.Expr(
		"COALESCE(tsv, to_tsvector(?, COALESCE(order_id,'') || ' ' || COALESCE(status,'') || ' ' || COALESCE(shipping_name,'') || ' ' || COALESCE(shipping_phone,'') || ' ' || COALESCE(shipping_address,'') || ' ' || COALESCE(shipping_zip_code,'')))",
		lang,
	)

	if err := s.db.Model(&orderRecord{}).
		Where("? @@ ?", tsvector, tsquery).
		Order(gorm.Expr("ts_rank(?, ?) DESC", tsvector, tsquery)).
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("search orders: %w", err)
	}

	results := make([]Result, len(rows))
	for i, r := range rows {
		results[i] = Result{
			ID:              r.ID,
			OrderID:         r.OrderID,
			UserID:          r.UserID,
			Status:          r.Status,
			ShippingName:    r.ShippingName,
			ShippingPhone:   r.ShippingPhone,
			ShippingAddress: r.ShippingAddress,
			ShippingZipCode: r.ShippingZipCode,
			TotalAmount:     r.TotalAmount,
		}
	}

	return results, nil
}
