package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	user "github.com/myproject/shop/internal/User"
	"github.com/myproject/shop/pkg/logger"
	"github.com/myproject/shop/pkg/utils"
	redis "github.com/redis/go-redis/v9"
)

const (
	CtxUserIDKey   = "userID"
	CtxUserRoleKey = "userRole"
	CtxTokenJTIKey = "tokenJTI"
)

var (
	JWTSecret   []byte
	RedisClient *redis.Client // 在应用启动时初始化
)

func InitJWT(secret string) {
	JWTSecret = []byte(secret)
}

func InitRedis(client *redis.Client) {
	RedisClient = client
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			logger.Warn("missing_authorization_header", map[string]interface{}{
				"method":  c.Request.Method,
				"path":    c.Request.URL.Path,
				"headers": c.Request.Header,
			})
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}
		tokenString := parts[1]
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		//超时检验
		if expVal, ok := claims["exp"]; ok {
			switch exp := expVal.(type) {
			case float64:
				if int64(exp) < time.Now().Unix() {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
					return
				}
			case int64:
				if exp < time.Now().Unix() {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
					return
				}
			case json.Number:
				if expInt, err := exp.Int64(); err == nil {
					if expInt < time.Now().Unix() {
						c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
						return
					}
				}
			default:
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid exp claim in token"})
				return
			}
		}

		jti, _ := claims["jti"].(string)
		var userID uint
		var userIDOK bool
		switch v := claims["sub"].(type) {
		case float64:
			userID = uint(v)
			userIDOK = true
		case int64:
			userID = uint(v)
			userIDOK = true
		case int:
			userID = uint(v)
			userIDOK = true
		case uint:
			userID = v
			userIDOK = true
		case string:
			if parsed, err := strconv.ParseUint(v, 10, 64); err == nil {
				userID = uint(parsed)
				userIDOK = true
			}
		case json.Number:
			if parsed, err := v.Int64(); err == nil {
				userID = uint(parsed)
				userIDOK = true
			}
		}
		if !userIDOK {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid sub claim in token"})
			return
		}
		var role uint
		var roleOK bool
		switch v := claims["role"].(type) {
		case float64:
			role = uint(v)
			roleOK = true
		case int64:
			role = uint(v)
			roleOK = true
		case int:
			role = uint(v)
			roleOK = true
		case uint:
			role = v
			roleOK = true
		case json.Number:
			if parsed, err := v.Int64(); err == nil {
				role = uint(parsed)
				roleOK = true
			}
		}
		if !roleOK {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid role claim in token"})
			return
		}
		// 检查黑名单
		if RedisClient != nil {
			ctx := context.Background()
			if exists, _ := RedisClient.Get(ctx, "bl:access:"+jti).Result(); exists != "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
				return
			}
		}

		c.Set(CtxUserIDKey, userID)
		c.Set(CtxUserRoleKey, role)
		c.Set(CtxTokenJTIKey, jti)
		c.Next()
	}
}

func MerchantAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roleVal, exists := ctx.Get(CtxUserRoleKey)
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permission denied: Role not found"})
			return
		}
		role, ok := roleVal.(uint)
		if !ok || (role != user.RoleMerchant && role != user.RoleAdmin) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permission denied: Only merchants and admins can perform this action"})
			return
		}
		ctx.Next()
	}
}
