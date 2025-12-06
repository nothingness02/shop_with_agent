package shop

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/myproject/shop/pkg/middleware"
)

const (
	shopTTL          = 5 * time.Minute
	productListTTL   = 5 * time.Minute
	productSingleTTL = 5 * time.Minute
)

type ShopService struct {
	rep   *ShopRepository
	cache *middleware.RedisStore
}

func NewShopService(rep *ShopRepository, cache *middleware.RedisStore) *ShopService {
	return &ShopService{rep: rep, cache: cache}
}

func (s *ShopService) List(limit, offset int) ([]Shop, error) {
	if limit <= 0 {
		limit = 20
	}
	cacheKey := fmt.Sprintf("shops:list:%d:%d", limit, offset)
	if s.cache != nil {
		var cached []Shop
		if ok, err := s.cache.GetObject(context.Background(), cacheKey, &cached); err == nil && ok {
			return cached, nil
		}
	}

	shops, err := s.rep.List(limit, offset)
	if err != nil || s.cache == nil {
		return shops, err
	}
	_ = s.cache.SetObjectWithTTL(context.Background(), cacheKey, shops, shopTTL)
	return shops, err
}

func (s *ShopService) GetShopByID(id uint) (*Shop, error) {
	var shop Shop
	cacheKey := "shop_" + strconv.FormatUint(uint64(id), 10)
	if s.cache != nil {
		if ok, err := s.cache.GetObject(context.Background(), cacheKey, &shop); err == nil && ok {
			return &shop, nil
		}
	}

	res, err := s.rep.Get(id)
	if err != nil || s.cache == nil {
		return res, err
	}
	_ = s.cache.SetObjectWithTTL(context.Background(), cacheKey, res, shopTTL)
	return res, nil
}

func (s *ShopService) CreateShop(sh *Shop) error {
	if sh == nil {
		return errors.New("shop is nil")
	}
	if sh.Name == "" {
		return errors.New("shop name is required")
	}
	if err := s.rep.Create(sh); err != nil {
		return err
	}
	if s.cache != nil {
		_ = s.cache.SetObjectWithTTL(context.Background(), "shop_"+strconv.FormatUint(uint64(sh.ID), 10), sh, shopTTL)
		_ = s.cache.DelteKey(context.Background(), "shops:list:20:0")
	}
	return nil
}

func (s *ShopService) UpdateShop(sh *Shop) error {
	if sh == nil || sh.ID == 0 {
		return errors.New("invalid shop")
	}
	if err := s.rep.Update(sh); err != nil {
		return err
	}
	if s.cache != nil {
		key := "shop_" + strconv.FormatUint(uint64(sh.ID), 10)
		_ = s.cache.SetObjectWithTTL(context.Background(), key, sh, shopTTL)
		_ = s.cache.DelteKey(context.Background(), "shops:list:20:0")
	}
	return nil
}

func (s *ShopService) Delete(id uint) error {
	if err := s.rep.Delete(id); err != nil {
		return err
	}
	if s.cache != nil {
		_ = s.cache.DelteKey(context.Background(), "shop_"+strconv.FormatUint(uint64(id), 10))
		_ = s.cache.DelteKey(context.Background(), "shops:list:20:0")
		_ = s.cache.DelteKey(context.Background(), fmt.Sprintf("shop_products:%d:20:0", id))
	}
	return nil
}

func (s *ShopService) BatchDelete(ids []uint) error {
	if err := s.rep.BatchDelete(ids); err != nil {
		return err
	}
	if s.cache != nil {
		for _, id := range ids {
			_ = s.cache.DelteKey(context.Background(), "shop_"+strconv.FormatUint(uint64(id), 10))
		}
		_ = s.cache.DelteKey(context.Background(), "shops:list:20:0")
	}
	return nil
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
	if err := s.rep.CreateProduct(shopID, p); err != nil {
		return err
	}
	if s.cache != nil {
		_ = s.cache.DelteKey(context.Background(), fmt.Sprintf("shop_products:%d:20:0", shopID))
	}
	return nil
}

func (s *ShopService) GetProductByCode(code uint) (*Product, error) {
	var product Product
	cacheKey := fmt.Sprintf("product:%d", code)
	if s.cache != nil {
		if ok, err := s.cache.GetObject(context.Background(), cacheKey, &product); err == nil && ok {
			return &product, nil
		}
	}

	res, err := s.rep.GetProductByCode(code)
	if err != nil || s.cache == nil {
		return res, err
	}
	_ = s.cache.SetObjectWithTTL(context.Background(), cacheKey, res, productSingleTTL)
	return res, nil
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
	cacheKey := fmt.Sprintf("shop_products:%d:%d:%d", shopID, limit, offset)
	if s.cache != nil {
		var cached []Product
		if ok, err := s.cache.GetObject(context.Background(), cacheKey, &cached); err == nil && ok {
			return cached, nil
		}
	}

	products, err := s.rep.ListProductsByShop(shopID, limit, offset)
	if err != nil || s.cache == nil {
		return products, err
	}
	_ = s.cache.SetObjectWithTTL(context.Background(), cacheKey, products, productListTTL)
	return products, err
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
	if err := s.rep.UpdateProduct(p); err != nil {
		return err
	}
	if s.cache != nil {
		_ = s.cache.SetObjectWithTTL(context.Background(), fmt.Sprintf("product:%d", p.ID), p, productSingleTTL)
		if p.ShopID != 0 {
			_ = s.cache.DelteKey(context.Background(), fmt.Sprintf("shop_products:%d:20:0", p.ShopID))
		}
	}
	return nil
}

func (s *ShopService) DeleteProduct(id uint) error {
	var p *Product
	if s.cache != nil || id > 0 {
		// best-effort fetch to know ShopID for invalidation; ignore error
		p, _ = s.rep.GetProductByCode(id)
	}
	if err := s.rep.DeleteProduct(id); err != nil {
		return err
	}
	if s.cache != nil {
		_ = s.cache.DelteKey(context.Background(), fmt.Sprintf("product:%d", id))
		if p != nil && p.ShopID != 0 {
			_ = s.cache.DelteKey(context.Background(), fmt.Sprintf("shop_products:%d:20:0", p.ShopID))
		}
	}
	return nil
}

func (s *ShopService) BatchDeleteProducts(ids []uint) error {
	if err := s.rep.BatchDeleteProducts(ids); err != nil {
		return err
	}
	if s.cache != nil {
		for _, id := range ids {
			_ = s.cache.DelteKey(context.Background(), fmt.Sprintf("product:%d", id))
		}
	}
	return nil
}
