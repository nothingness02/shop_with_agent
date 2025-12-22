package Coordinator

import (
	"context"

	"github.com/myproject/shop/internal/Order"
	shop "github.com/myproject/shop/internal/Shop"
	"gorm.io/gorm"
)

type CheckoutService struct {
	db           *gorm.DB
	orderService *Order.OrderService
	shopService  *shop.ShopService
}

func NewCheckoutService(db *gorm.DB, orderS *Order.OrderService, shopS *shop.ShopService) *CheckoutService {
	return &CheckoutService{
		db:           db,
		orderService: orderS,
		shopService:  shopS,
	}
}

func (s *CheckoutService) PlaceOrder(ctx context.Context, order *Order.Order) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		//扣除库存
		for _, item := range order.OrderItems {
			if err := s.shopService.DecreaseStockWithTx(ctx, tx, item.ProductID, item.Quantity); err != nil {
				return err
			}
		}
		//创建订单
		if err := s.orderService.CreateOrderWithTx(ctx, tx, order); err != nil {
			return err // 触发事务回滚
		}
		return nil
	})
}
