package Coordinator

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myproject/shop/internal/Order"
)

type TradeHandler struct {
	service *CheckoutService
}

type createOrderRequest struct {
	UserID uint              `json:"user_id" binding:"required"`
	Items  []createOrderItem `json:"items" binding:"required,dive"`
}

type createOrderItem struct {
	ProductID   uint    `json:"product_id" binding:"required"`
	ProductName string  `json:"product_name" binding:"required"`
	ProductImg  string  `json:"product_img"`
	Price       float64 `json:"price" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required,min=1"`
}

func NewTradeHandler(srv *CheckoutService) *TradeHandler {
	return &TradeHandler{service: srv}
}

func (h *TradeHandler) CreateOrder(c *gin.Context) {
	var request createOrderRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	o := Order.Order{
		UserID:     request.UserID,
		Status:     Order.OrderStatusPending,
		OrderItems: make([]Order.OrderItem, len(request.Items)),
	}
	for i := range request.Items {
		o.OrderItems[i] = Order.OrderItem{
			ProductID:   request.Items[i].ProductID,
			ProductName: request.Items[i].ProductName,
			ProductImg:  request.Items[i].ProductImg,
			Price:       request.Items[i].Price,
			Quantity:    request.Items[i].Quantity,
			Subtotal:    request.Items[i].Price * float64(request.Items[i].Quantity),
		}
	}
	if err := h.service.PlaceOrder(context.Background(), &o); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, o)
}
