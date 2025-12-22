package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	auth "github.com/myproject/shop/internal/Auth"
	cart "github.com/myproject/shop/internal/Cart"
	comment "github.com/myproject/shop/internal/Comment"
	Coordinator "github.com/myproject/shop/internal/Coordinator"
	"github.com/myproject/shop/internal/Order"
	shop "github.com/myproject/shop/internal/Shop"
	user "github.com/myproject/shop/internal/User"
	config "github.com/myproject/shop/internal/config"
	search "github.com/myproject/shop/internal/search"
	"github.com/myproject/shop/pkg/logger"
	"github.com/myproject/shop/pkg/middleware"
)

type Application struct {
	*gin.Engine
	config *config.Config
}

func NewApplication(cfg *config.Config,
	userH *user.UserHandle,
	authH *auth.AuthHandler,
	orderH *Order.OrderHandler,
	shopH *shop.ShopHandler,
	searchH *search.Handler,
	commentH *comment.CommentHandler,
	cartH *cart.CartHandler,
	coordinatorH *Coordinator.TradeHandler) *Application {
	gin.SetMode(cfg.Server.Mode)
	app := &Application{
		Engine: gin.New(),
		config: cfg,
	}
	app.Use(gin.Recovery())
	app.Use(logger.GinLogger())
	app.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	middleware.InitJWT(app.config.JWT.Secret)
	v0 := app.Group("/api/v0")
	{
		// User routes
		v0.POST("/users", userH.RegisterUser)
		v0.GET("/users/:id", userH.GetUserByID)
		v0.PATCH("/users/:id", userH.UpdateUser)
		v0.DELETE("/users/:id", userH.DeleteUser)
		// v0.GET("/users/:username", userHandler.GetUserByName)
		// Auth routes
		v0.POST("/auth/login", authH.Login)
		v0.POST("/auth/refresh", authH.Refresh)
		v0.POST("/auth/logout", middleware.JWTAuthMiddleware(), authH.Logout)
	}

	v1 := app.Group("/api/v1")
	v1.Use(middleware.JWTAuthMiddleware())
	{
		v1.GET("/orders", orderH.ListOrders)
		v1.GET("/orders/:id", coordinatorH.CreateOrder)
		v1.POST("/orders", orderH.CreateOrder)
		v1.PATCH("/orders/:id/status", orderH.UpdateOrderStatus)
		v1.DELETE("/orders/:id", orderH.DeleteOrder)

		// Cart routes
		v1.GET("/cart", cartH.List)
		v1.POST("/cart/items", cartH.Add)
		v1.PATCH("/cart/items/:id", cartH.Update)
		v1.DELETE("/cart/items/:id", cartH.Delete)
		v1.DELETE("/cart", cartH.Clear)
	}

	v2 := app.Group("/api/v2")
	v2.Use(middleware.JWTAuthMiddleware())

	// Customer-facing (authenticated) routes
	v2Customer := v2.Group("")
	{
		// Search routes
		v2Customer.GET("/search/products", searchH.SearchProducts)
		v2Customer.GET("/search/orders", searchH.SearchOrders)

		// Shop/product discovery
		v2Customer.GET("/shops", shopH.ListShops)
		v2Customer.GET("/shops/:id/products", shopH.ListProducts)
		v2Customer.GET("/shops/:id/products/search", shopH.GetProductByName)
		v2Customer.GET("/products/:id", shopH.GetProductByCode)
		v2Customer.GET("/shops/:id", shopH.GetShop)
	}

	// Merchant/admin routes (require role check) – paths unchanged
	v2Merchant := v2.Group("")
	v2Merchant.Use(middleware.MerchantAuthMiddleware())
	{
		// Shop & product management
		v2Merchant.POST("/shops", shopH.CreateShop)
		v2Merchant.POST("/shops/:id/products", shopH.CreateProduct)
		v2Merchant.PATCH("/products/:id", shopH.UpdateProduct)
		v2Merchant.DELETE("/products/:id", shopH.DeleteProduct)
		v2Merchant.DELETE("/products", shopH.BatchDeleteProducts)
		v2Merchant.PATCH("/shops/:id", shopH.UpdateShop)
		v2Merchant.DELETE("/shops/:id", shopH.DeleteShop)
		v2Merchant.DELETE("/shops", shopH.BatchDeleteShops)
	}

	v3 := app.Group("api/v3")
	{
		v3.GET("/shops/:id/comments", commentH.ListCommentsByShop)
		v3.POST("/shops/:id/comments", middleware.JWTAuthMiddleware(), commentH.CreateComment)
		v3.DELETE("/shops/:id/comments", commentH.DeleteComments)
	}

	return app
}

func (app *Application) run() error {
	err := app.Run(app.config.Server.RunAddr)
	if err != nil {
		return err
	}
	return nil
}
