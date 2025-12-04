package shop

import (
	"errors"
)

type ShopService struct {
	rep *ShopRepository
}

func NewShopService(rep *ShopRepository) *ShopService {
	return &ShopService{rep: rep}
}

func (s *ShopService) List(limit, offset int) ([]Shop, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.rep.List(limit, offset)
}

func (s *ShopService) GetShopByID(id uint) (*Shop, error) {
	return s.rep.Get(id)
}

func (s *ShopService) CreateShop(sh *Shop) error {
	if sh == nil {
		return errors.New("shop is nil")
	}
	if sh.Name == "" {
		return errors.New("shop name is required")
	}
	return s.rep.Create(sh)
}

func (s *ShopService) UpdateShop(sh *Shop) error {
	if sh == nil || sh.ID == 0 {
		return errors.New("invalid shop")
	}
	return s.rep.Update(sh)
}

func (s *ShopService) Delete(id uint) error {
	return s.rep.Delete(id)
}

func (s *ShopService) BatchDelete(ids []uint) error {
	return s.rep.BatchDelete(ids)
}

// Product-related methods
func (s *ShopService) CreateProduct(shopID uint, p *Product) error {
	if shopID == 0 {
		return errors.New("shop id is required")
	}
	if p == nil {
		return errors.New("product is nil")
	}
	if p.Name == "" {
		return errors.New("product name is required")
	}
	if p.Price <= 0 {
		return errors.New("product price must be greater than 0")
	}
	return s.rep.CreateProduct(shopID, p)
}

func (s *ShopService) GetProductByCode(code uint) (*Product, error) {
	return s.rep.GetProductByCode(code)
}

func (s *ShopService) GetProductByName(shopID uint, name string) (*Product, error) {
	if shopID == 0 || name == "" {
		return nil, errors.New("shop id and product name are required")
	}
	return s.rep.GetProductByName(shopID, name)
}

func (s *ShopService) ListProductsByShop(shopID uint, limit, offset int) ([]Product, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.rep.ListProductsByShop(shopID, limit, offset)
}

func (s *ShopService) UpdateProduct(p *Product) error {
	if p == nil || p.ID == 0 {
		return errors.New("invalid product")
	}
	if p.Name == "" {
		return errors.New("product name is required")
	}
	if p.Price <= 0 {
		return errors.New("product price must be greater than 0")
	}
	return s.rep.UpdateProduct(p)
}

func (s *ShopService) DeleteProduct(id uint) error {
	return s.rep.DeleteProduct(id)
}

func (s *ShopService) BatchDeleteProducts(ids []uint) error {
	return s.rep.BatchDeleteProducts(ids)
}
