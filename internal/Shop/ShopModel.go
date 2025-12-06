package shop

import (
	"gorm.io/gorm"
)

type Shop struct {
	gorm.Model
	Name        string    `gorm:"size:100;not null"` // 商店名称Name
	Description string    `gorm:"size:255"`          // 商店描述
	OwnerID     uint      `gorm:"index"`             // 店主用户ID
	Products    []Product `gorm:"foreignKey:ShopID"` // 关联的商品
}

type Product struct {
	gorm.Model
	ShopID uint `gorm:"index"` // 所属商店ID
	// ProductID   uint    `gorm:"uniqueIndex;size:100"` // 商品ID
	Name        string  `gorm:"size:100;not null"`  // 商品名称
	Description string  `gorm:"size:255"`           // 商品描述
	Price       float64 `gorm:"type:decimal(10,2)"` // 商品价格
	Stock       int     // 库存数量
	ProductImg  string  `gorm:"size:500"`                         // 商品图片URL
	Tsv         string  `gorm:"type:tsvector;index:,type:gin;->"` // 用于全文搜索, GORM不会写入，由数据库触发器填充
}

type Category struct {
	gorm.Model
	Name        string    `gorm:"size:100;not null"`             // 分类名称
	Description string    `gorm:"size:255"`                      // 分类描述
	Products    []Product `gorm:"many2many:product_categories;"` // 关联的商品
}
