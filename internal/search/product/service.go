package product

import (
	"fmt"

	"gorm.io/gorm"
)

// Service handles product full-text search.
type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// Result is a lightweight view returned to callers.
type Result struct {
	ID          uint    `json:"id"`
	ShopID      uint    `json:"shop_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	ProductImg  string  `json:"product_img"`
}

// productRecord mirrors the database schema for search queries.
type productRecord struct {
	gorm.Model
	ShopID      uint
	Name        string
	Description string
	Price       float64
	Stock       int
	ProductImg  string
	Tsv         string `gorm:"type:tsvector;index:,type:gin;->"`
}

func (productRecord) TableName() string {
	return "products"
}

// Search runs a full-text search against products.
func (s *Service) Search(query, lang string) ([]Result, error) {
	var rows []productRecord

	tsquery := gorm.Expr("plainto_tsquery(?, ?)", lang, query)

	if err := s.db.Model(&productRecord{}).
		Where("tsv @@ ?", tsquery).
		Order(gorm.Expr("ts_rank(tsv, ?) DESC", tsquery)).
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("search products: %w", err)
	}

	results := make([]Result, len(rows))
	for i, r := range rows {
		results[i] = Result{
			ID:          r.ID,
			ShopID:      r.ShopID,
			Name:        r.Name,
			Description: r.Description,
			Price:       r.Price,
			Stock:       r.Stock,
			ProductImg:  r.ProductImg,
		}
	}

	return results, nil
}
