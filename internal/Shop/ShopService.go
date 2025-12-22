package shop

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/myproject/shop/pkg/middleware"
	"gorm.io/gorm"
)

const (
	shopTTL          = 5 * time.Minute
	productListTTL   = 5 * time.Minute
	productSingleTTL = 5 * time.Minute
)

// 返回值:
//
//	-1: 库存 Key 不存在 (需要预热)
//	-2: 库存不足
//	>=0: 扣减后的剩余库存
const luaScript = `
local stock = redis.call("get", KEYS[1])
if not stock then
    return -1
end
stock = tonumber(stock)
local qty = tonumber(ARGV[1])
if stock < qty then
	return -2
end
return redis.call("decrby",KEYS[1],qty)
`

type ShopService struct {
	rep   *ShopRepository
	cache *middleware.RedisStore
	sf    singleflight.Group
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

//===================Product===================================================

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
		stockKey := fmt.Sprintf("product:stock:%d", p.ID)
		s.cache.SetObjectWithTTL(context.Background(), stockKey, p.Stock, 0) // 0 表示永不过期
		_ = s.cache.DelteKey(context.Background(), fmt.Sprintf("shop_products:%d:20:0", shopID))
	}
	return nil
}

// 添加了旁路缓存机制
func (s *ShopService) GetProductByCode(code uint) (*Product, error) {
	var product Product
	cacheKey := fmt.Sprintf("product:%d", code)
	//先查询缓存是否存在
	if s.cache != nil {
		if ok, err := s.cache.GetObject(context.Background(), cacheKey, &product); err == nil && ok {
			return &product, nil
		}
	}
	//如果不存在，合并请求流为数据库减压（slow_path）
	result, err, _ := s.sf.Do(cacheKey, func() (any, error) {
		//查询一次数据库
		res, err := s.rep.GetProductByCode(code)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.cache.SetObjectWithTTL(context.Background(), cacheKey, "{}", 5*time.Minute)
				return nil, err
			}
			return nil, err
		}
		// 4. 写入缓存 (Cache-Aside)
		s.cache.SetObjectWithTTL(context.Background(), cacheKey, res, 1*time.Hour)
		return res, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*Product), err
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
		//清除失效缓存
		if p.ShopID != 0 {
			_ = s.cache.DelteKey(context.Background(), fmt.Sprintf("shop_products:%d:20:0", p.ShopID))
		}

		stockKey := fmt.Sprintf("product:stock:%d", p.ID)
		s.cache.Client.Set(context.Background(), stockKey, p.Stock, 0)
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

func (s *ShopService) DecreaseStockWithTx(ctx context.Context, tx *gorm.DB, productid uint, quantity int) error {
	//构架cache key
	cachekey := fmt.Sprintf("product:stock:%d", productid)
	if s.cache != nil {
		//第一次防护先将库存减在redis里面
		// keys: [cacheKey], args: [quantity]
		res, err := s.cache.Client.Eval(context.Background(), luaScript, []string{cachekey}, quantity).Result()
		if err != nil {
			fmt.Printf("Redis error: %v, falling back to DB\n", err)
		} else {
			retVal := res.(int64)
			switch retVal {
			case -1:
				//Key 不存在：说明缓存过期或者没有进行数据预热
			case -2:
				return errors.New("insufficient stock (redis)")
			default:
				//获取扣除数据库的资格
			}
		}
	}

	//数据库阶段
	err := s.rep.DecreaseStockWithTx(ctx, tx, productid, quantity)
	if err != nil {
		//4.数据库扣除失败必须将缓存中扣除的库存返还回去
		if s.cache != nil {
			s.cache.Client.IncrBy(context.Background(), cachekey, int64(quantity))
		}
		return fmt.Errorf("decrease stock failed: %w", err)
	}
	// 扣减成功，可以清理商品详情缓存，保证数据新鲜度
	if s.cache != nil {
		_ = s.cache.DelteKey(context.Background(), fmt.Sprintf("product:%d", productid))
	}
	return nil
}
