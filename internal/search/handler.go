package search

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ordersearch "github.com/myproject/shop/internal/search/order"
	"github.com/myproject/shop/internal/search/product"
)

// Handler exposes HTTP endpoints for search.
type Handler struct {
	productService *product.Service
	orderService   *ordersearch.Service
}

func NewHandler(productService *product.Service, orderService *ordersearch.Service) *Handler {
	return &Handler{
		productService: productService,
		orderService:   orderService,
	}
}

func (h *Handler) SearchProducts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}
	lang := c.DefaultQuery("lang", "simple")

	results, err := h.productService.Search(query, lang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search products"})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (h *Handler) SearchOrders(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}
	lang := c.DefaultQuery("lang", "simple")

	results, err := h.orderService.Search(query, lang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search orders"})
		return
	}

	c.JSON(http.StatusOK, results)
}
