package Order

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service *OrderService
}

func NewOrderHandler(service *OrderService) *OrderHandler {
	return &OrderHandler{service: service}
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

type updateStatusReq struct {
	Status OrderStatus `json:"status" binding:"required"`
}

type batchDeleteReq struct {
	IDs []uint `json:"ids" binding:"required,min=1,dive"`
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	oreders, err := h.service.List(50, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, oreders)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	o, err := h.service.GetOrderById(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, o)
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	o := Order{
		UserID:     req.UserID,
		Status:     OrderStatusPending,
		OrderItems: make([]OrderItem, len(req.Items)),
	}
	for i := range req.Items {
		o.OrderItems[i] = OrderItem{
			ProductID:   req.Items[i].ProductID,
			ProductName: req.Items[i].ProductName,
			ProductImg:  req.Items[i].ProductImg,
			Price:       req.Items[i].Price,
			Quantity:    req.Items[i].Quantity,
			Subtotal:    req.Items[i].Price * float64(req.Items[i].Quantity),
		}
	}
	if err := h.service.CreateOrder(&o); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, o)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req updateStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateStatus(uint(id), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order deleted"})
}

func (h *OrderHandler) BatchDeleteOrders(c *gin.Context) {
	var req batchDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.BatchDelete(req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "orders deleted"})
}
