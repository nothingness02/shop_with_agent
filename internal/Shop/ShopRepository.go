package shop

import (
	"errors"

	"github.com/myproject/shop/pkg/database"
	"gorm.io/gorm"
)

type ShopRepository struct {
	Database *database.Database
}

func NewRepository(db *database.Database) *ShopRepository {
	return &ShopRepository{Database: db}
}

func (r *ShopRepository) List(limit, offset int) ([]Shop, error) {
	var shops []Shop
	if err := r.Database.DB.Preload("Products").Limit(limit).Offset(offset).Find(&shops).Error; err != nil {
		return nil, err
	}
	return shops, nil
}

func (r *ShopRepository) Get(id uint) (*Shop, error) {
	var s Shop
	if err := r.Database.DB.Preload("Products").First(&s, id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *ShopRepository) Create(s *Shop) error {
	return r.Database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(s).Error; err != nil {
			return err
		}
		for i := range s.Products {
			s.Products[i].ShopID = s.ID
		}
		if len(s.Products) > 0 {
			if err := tx.Create(&s.Products).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *ShopRepository) Update(s *Shop) error {

	return r.Database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(s).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *ShopRepository) Delete(id uint) error {
	return r.Database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("shop_id = ?", id).Delete(&Product{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&Shop{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *ShopRepository) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("empty ids")
	}
	return r.Database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("shop_id IN ?", ids).Delete(&Product{}).Error; err != nil {
			return err
		}
		return tx.Delete(&Shop{}, ids).Error
	})
}

// Product-related methods
func (r *ShopRepository) CreateProduct(shopID uint, p *Product) error {
	p.ShopID = shopID
	return r.Database.DB.Create(p).Error
}

func (r *ShopRepository) GetProductByCode(code uint) (*Product, error) {
	var p Product
	if err := r.Database.DB.First(&p, "id = ?", code).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ShopRepository) GetProductByName(shopID uint, name string) (*Product, error) {
	var p Product
	if err := r.Database.DB.Where("shop_id = ? AND name = ?", shopID, name).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ShopRepository) ListProductsByShop(shopID uint, limit, offset int) ([]Product, error) {
	var products []Product
	if err := r.Database.DB.Where("shop_id = ?", shopID).Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ShopRepository) UpdateProduct(p *Product) error {
	if p.ShopID == 0 || p == nil {
		return errors.New("invalid product")
	}
	return r.Database.DB.Model(p).Select("name", "description", "price", "stock", "product_img").Save(p).Error
}

func (r *ShopRepository) DeleteProduct(id uint) error {
	return r.Database.DB.Delete(&Product{}, id).Error
}

func (r *ShopRepository) BatchDeleteProducts(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("empty ids")
	}
	return r.Database.DB.Delete(&Product{}, ids).Error
}

// 数据库层面的乐观锁扣减
func (r *ShopRepository) DecreaseStock(id uint, quannity int) error {
	if quannity <= 0 {
		return errors.New("quantity must be greater than 0")
	}
	result := r.Database.DB.Model(&Product{}).
		Where("id = ? AND stock >= ?", id, quannity).
		Update("stock", gorm.Expr("stock - ?", quannity))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("insufficient stock")
	}
	return nil
}

// 回滚库存
func (r *ShopRepository) AddStock(id uint, quannity int) error {
	return r.Database.DB.Model(&Product{}).Where("id = ?", id).Update("stock", gorm.Expr("stock + ?", quannity)).Error

}
