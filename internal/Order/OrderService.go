package Order

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type OrderService struct {
	rep *OrderRepository
}

func NewOrderService(rep *OrderRepository) *OrderService {
	return &OrderService{rep: rep}
}
func (s *OrderService) List(limit, offset int) ([]Order, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.rep.List(limit, offset)
}
func (s *OrderService) GetOrderById(id uint) (*Order, error) {
	return s.rep.Get(id)
}

func (s *OrderService) CreateOrder(o *Order) error {
	if o == nil || len(o.OrderItems) == 0 {

		return errors.New("order items is empty")
	}
	return s.rep.Create(o)

}

func (s *OrderService) CreateOrderWithTx(ctx context.Context, tx *gorm.DB, o *Order) error {
	if o == nil || len(o.OrderItems) == 0 {

		return errors.New("order items is empty")
	}
	return s.rep.CreateWithTx(ctx, tx, o)
}

func (s *OrderService) UpdateStatus(id uint, status OrderStatus) error {
	return s.rep.UpdateStatus(id, status)
}

func (s *OrderService) Delete(id uint) error {
	return s.rep.Delete(id)
}

func (s *OrderService) BatchDelete(ids []uint) error {
	return s.rep.BatchDelete(ids)
}
