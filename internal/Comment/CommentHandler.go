package comment

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/myproject/shop/pkg/middleware"
)

type CommentHandler struct {
	service *CommentService
}

func NewCommentHandler(srv *CommentService) *CommentHandler {
	return &CommentHandler{service: srv}
}

type createCommentReq struct {
	Content         string `json:"content" binding:"required"`
	ParentCommentID *uint  `json:"parent_comment_id"`
}

// ListCommentsByShop GET /shops/:id/comments
// 返回格式：{ "total": N, "items": [...] }
func (h *CommentHandler) ListCommentsByShop(c *gin.Context) {
	shopIDInt, _ := strconv.Atoi(c.Param("id"))
	shopID := uint(shopIDInt)

	var parentID *uint
	if p := c.Query("parent_id"); p != "" {
		if pid, err := strconv.ParseUint(p, 10, 64); err == nil {
			tmp := uint(pid)
			parentID = &tmp
		}
	}

	page := 1
	pageSize := 20
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	}

	items, total, err := h.service.ListCommentsByShopWithCount(context.Background(), shopID, parentID, page, pageSize, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": total, "items": items})
}

// CreateComment POST /shops/:id/comments
func (h *CommentHandler) CreateComment(c *gin.Context) {
	shopIDInt, _ := strconv.Atoi(c.Param("id"))
	shopID := uint(shopIDInt)

	var req createCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从中间件获取 user id（string 存储），尝试解析
	var userID uint
	if v, ok := c.Get(middleware.CtxUserIDKey); ok {
		if s, ok := v.(string); ok {
			if uid, err := strconv.ParseUint(s, 10, 64); err == nil {
				userID = uint(uid)
			}
		}
	}

	cm := &Comment{
		UserID:          userID,
		ShopID:          shopID,
		Content:         req.Content,
		ParentCommentID: req.ParentCommentID,
	}
	if err := h.service.CreateComment(context.Background(), cm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cm)
}

// DeleteComments DELETE /shops/:id/comments
// Query param: cascade=true|false
func (h *CommentHandler) DeleteComments(c *gin.Context) {
	shopIDInt, _ := strconv.Atoi(c.Param("id"))
	shopID := uint(shopIDInt)

	// 从中间件获取 user id（string 存储），尝试解析
	var userID uint
	if v, ok := c.Get(middleware.CtxUserIDKey); ok {
		if s, ok := v.(string); ok {
			if uid, err := strconv.ParseUint(s, 10, 64); err == nil {
				userID = uint(uid)
			}
		}
	}

	cascade := false
	if c.Query("cascade") == "true" {
		cascade = true
	}

	if err := h.service.DeleteComment(context.Background(), userID, shopID, cascade); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "comments deleted"})
}