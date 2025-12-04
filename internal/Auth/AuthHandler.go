package auth

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/myproject/shop/pkg/middleware"
)

type AuthHandler struct {
	svc *AuthService
}

func NewAuthHandler(svc *AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	access, refresh, userID, role, err := h.svc.Login(context.Background(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": access, "refresh_token": refresh, "user_id": userID, "role": role})
}

type refreshReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	access, refresh, _, _, err := h.svc.Refresh(context.Background(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": access, "refresh_token": refresh})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// require JWTAuthMiddleware to set context
	userID := uint(0)
	if v, ok := c.Get(middleware.CtxUserIDKey); ok {
		switch vv := v.(type) {
		case string:
			if parsed, err := strconv.ParseUint(vv, 10, 64); err == nil {
				userID = uint(parsed)
			}
		case float64:
			userID = uint(vv)
		case int:
			userID = uint(vv)
		case uint:
			userID = vv
		}
	}
	jti := ""
	if v, ok := c.Get(middleware.CtxTokenJTIKey); ok {
		if s, ok := v.(string); ok {
			jti = s
		}
	}
	_ = h.svc.Logout(context.Background(), userID, jti)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
