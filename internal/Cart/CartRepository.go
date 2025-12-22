package cart

import (
	"errors"

	shop "github.com/myproject/shop/internal/Shop"
	"github.com/myproject/shop/pkg/database"
	"gorm.io/gorm"
)

type CartRepository struct {
	Database *database.Database
}

func NewCartRepository(db *database.Database) *CartRepository {
	return &CartRepository{Database: db}
}

func (r *CartRepository) ListByUser(userID uint) ([]CartItem, error) {
	var items []CartItem
	if err := r.Database.DB.Where("user_id = ?", userID).Order("id desc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *CartRepository) AddOrUpdate(userID, productID uint, quantity int) (*CartItem, error) {
	if quantity <= 0 {
		return nil, errors.New("quantity must be positive")
	}
	var product shop.Product
	if err := r.Database.DB.First(&product, productID).Error; err != nil {
		return nil, err
	}

	var existing CartItem
	err := r.Database.DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&existing).Error
	if err == nil {
		existing.Quantity += quantity
		return &existing, r.Database.DB.Save(&existing).Error
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	item := CartItem{
		UserID:      userID,
		ProductID:   productID,
		ProductName: product.Name,
		ProductImg:  product.ProductImg,
		Price:       product.Price,
		Quantity:    quantity,
	}
	if err := r.Database.DB.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CartRepository) UpdateQuantity(userID, itemID uint, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	return r.Database.DB.Model(&CartItem{}).Where("id = ? AND user_id = ?", itemID, userID).Update("quantity", quantity).Error
}

func (r *CartRepository) DeleteItem(userID, itemID uint) error {
	return r.Database.DB.Where("id = ? AND user_id = ?", itemID, userID).Delete(&CartItem{}).Error
}

func (r *CartRepository) Clear(userID uint) error {
	return r.Database.DB.Where("user_id = ?", userID).Delete(&CartItem{}).Error
}