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
// Result keeps the same field names/json keys as the Product model
// so that front-end rendering (which consumes list products API) works.
type Result struct {
	ID          uint    `json:"ID"`
	ShopID      uint    `json:"ShopID"`
	Name        string  `json:"Name"`
	Description string  `json:"Description"`
	Price       float64 `json:"Price"`
	Stock       int     `json:"Stock"`
	ProductImg  string  `json:"ProductImg"`
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
	// 使用 coalesce(tsv, to_tsvector(...))，避免 tsv 未被触发器填充时无法命中
	tsvector := gorm.Expr("COALESCE(tsv, to_tsvector(?, COALESCE(name,'') || ' ' || COALESCE(description,'')))", lang)

	if err := s.db.Model(&productRecord{}).
		Where("? @@ ?", tsvector, tsquery).
		Order(gorm.Expr("ts_rank(?, ?) DESC", tsvector, tsquery)).
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
