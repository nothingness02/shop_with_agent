package Order

import (
	"errors"

	"github.com/myproject/shop/pkg/database"
	"gorm.io/gorm"
)

type OrderRepository struct {
	Database *database.Database
}

func NewRepository(db *database.Database) *OrderRepository {
	return &OrderRepository{Database: db}
}

func (r *OrderRepository) List(limit, offset int) ([]Order, error) {
	var orders []Order
	if err := r.Database.DB.Preload("Items").Limit(limit).Offset(offset).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) Get(id uint) (*Order, error) {
	var o Order
	if err := r.Database.DB.Preload("Items").First(&o, id).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OrderRepository) Create(o *Order) error {
	return r.Database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(o).Error; err != nil {
			return err
		}
		for i := range o.OrderItems {
			o.OrderItems[i].OrderID = o.ID
		}
		if len(o.OrderItems) > 0 {
			if err := tx.Create(o.OrderItems).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
func (r *OrderRepository) UpdateStatus(id uint, status OrderStatus) error {
	return r.Database.DB.Model(&Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *OrderRepository) Delete(id uint) error {
	return r.Database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("order_id = ?", id).Delete(&OrderItem{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&Order{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}
func (r *OrderRepository) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("empty ids")
	}
	return r.Database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("order_id IN ?", ids).Delete(&OrderItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(&Order{}, ids).Error
	})
}
