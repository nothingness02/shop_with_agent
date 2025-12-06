package shop

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShopHandler struct {
	service *ShopService
}

func NewShopHandler(service *ShopService) *ShopHandler {
	return &ShopHandler{service: service}
}

type createProductReq struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	Stock       int     `json:"stock"`
	ProductImg  string  `json:"product_img"`
}

type createShopReq struct {
	Name        string             `json:"name" binding:"required"`
	Description string             `json:"description"`
	OwnerID     uint               `json:"owner_id" binding:"required"`
	Products    []createProductReq `json:"products"`
}

type updateShopReq struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Products    []createProductReq `json:"products"`
}

type batchDeleteReq struct {
	IDs []uint `json:"ids" binding:"required,min=1,dive"`
}

func (h *ShopHandler) ListShops(c *gin.Context) {
	shops, err := h.service.List(50, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, shops)
}

func (h *ShopHandler) GetShop(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	s, err := h.service.GetShopByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *ShopHandler) CreateShop(c *gin.Context) {
	var req createShopReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sh := Shop{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     req.OwnerID,
		Products:    make([]Product, len(req.Products)),
	}
	for i := range req.Products {
		sh.Products[i] = Product{
			Name:        req.Products[i].Name,
			Description: req.Products[i].Description,
			Price:       req.Products[i].Price,
			Stock:       req.Products[i].Stock,
			ProductImg:  req.Products[i].ProductImg,
		}
	}
	if err := h.service.CreateShop(&sh); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sh)
}

func (h *ShopHandler) UpdateShop(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req updateShopReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sh := Shop{
		Model:       gorm.Model{ID: uint(id)},
		Name:        req.Name,
		Description: req.Description,
	}
	if len(req.Products) > 0 {
		sh.Products = make([]Product, len(req.Products))
		for i := range req.Products {
			sh.Products[i] = Product{
				Name:        req.Products[i].Name,
				Description: req.Products[i].Description,
				Price:       req.Products[i].Price,
				Stock:       req.Products[i].Stock,
				ProductImg:  req.Products[i].ProductImg,
			}
		}
	}
	if err := h.service.UpdateShop(&sh); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "shop updated"})
}

func (h *ShopHandler) DeleteShop(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "shop deleted"})
}

func (h *ShopHandler) BatchDeleteShops(c *gin.Context) {
	var req batchDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.BatchDelete(req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "shops deleted"})
}

// Product handlers - avoid redundant shop updates
func (h *ShopHandler) CreateProduct(c *gin.Context) {
	shopID, _ := strconv.Atoi(c.Param("id"))
	var req createProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ProductImg:  req.ProductImg,
	}
	if err := h.service.CreateProduct(uint(shopID), &p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *ShopHandler) GetProductByCode(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	p, err := h.service.GetProductByCode(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *ShopHandler) GetProductByName(c *gin.Context) {
	shopID, _ := strconv.Atoi(c.Param("id"))
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product name is required"})
		return
	}
	p, err := h.service.GetProductByName(uint(shopID), name)
	if err.Error() == "record not found" {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *ShopHandler) ListProducts(c *gin.Context) {
	shopID, _ := strconv.Atoi(c.Param("id"))
	products, err := h.service.ListProductsByShop(uint(shopID), 50, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

type updateProductReq struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	ProductImg  string  `json:"product_img"`
}

func (h *ShopHandler) UpdateProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req updateProductReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := Product{
		Model:       gorm.Model{ID: uint(id)},
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ProductImg:  req.ProductImg,
	}
	if err := h.service.UpdateProduct(&p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product updated"})
}

func (h *ShopHandler) DeleteProduct(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.DeleteProduct(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}

func (h *ShopHandler) BatchDeleteProducts(c *gin.Context) {
	var req batchDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.BatchDeleteProducts(req.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "products deleted"})
}
