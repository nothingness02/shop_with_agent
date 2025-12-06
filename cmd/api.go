package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	auth "github.com/myproject/shop/internal/Auth"
	comment "github.com/myproject/shop/internal/Comment"
	"github.com/myproject/shop/internal/Order"
	shop "github.com/myproject/shop/internal/Shop"
	user "github.com/myproject/shop/internal/User"
	search "github.com/myproject/shop/internal/search"
	ordersearch "github.com/myproject/shop/internal/search/order"
	"github.com/myproject/shop/internal/search/product"
	"github.com/myproject/shop/pkg/database"
	"github.com/myproject/shop/pkg/middleware"
)

type Application struct {
	*gin.Engine
	config *config
}
type config struct {
	address string
	db      *database.Dbconfig
}

func newConfig() config {
	return config{
		address: ":8080",
		db: &database.Dbconfig{
			Dsn: "postgres://postgres:postgres@localhost:5432/app_db?sslmode=disable",
		},
	}
}

func NewApplication(config *config) *Application {
	app := &Application{
		Engine: gin.Default(),
		config: config,
	}
	app.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	app.mountRoutes()
	return app
}

func (app *Application) mountRoutes() {

	db, err := database.NewDB(app.config.db)

	if err != nil {
		log.Fatal(err)
		return
	}

	if err := db.AutoMigrate(&Order.Order{}, &Order.OrderItem{},
		&shop.Shop{}, &shop.Product{},
		&shop.Category{}, &user.User{}, &comment.Comment{},
	); err != nil {
		log.Fatal(err)
		return
	}

	redisStore := middleware.NewRedisStore("localhost:6379", "redis", 0)
	middleware.InitRedis(redisStore.Client)

	productSearchService := product.NewService(db.DB)
	orderSearchService := ordersearch.NewService(db.DB)
	searchHandler := search.NewHandler(productSearchService, orderSearchService)
	orderrep := Order.NewRepository(db)
	orderService := Order.NewOrderService(orderrep)
	orderHandler := Order.NewOrderHandler(orderService)
	shoprep := shop.NewRepository(db)
	shopService := shop.NewShopService(shoprep, redisStore)
	shopHandler := shop.NewShopHandler(shopService)

	// 初始化用户仓库和 Redis（用于 refresh token 存储 / blacklist）
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewUserHandle(userService)

	// 初始化 Auth 服务与处理器
	authService := auth.NewAuthService(userRepo, redisStore)
	authHandler := auth.NewAuthHandler(authService)

	//初始化评论服务
	commentRepo := comment.NewRepository(db)
	commentService := comment.NewCommentService(commentRepo)
	commentHandler := comment.NewCommentHandler(commentService)

	v0 := app.Group("/api/v0")
	{
		// User routes
		v0.POST("/users", userHandler.RegisterUser)
		v0.GET("/users/:id", userHandler.GetUserByID)
		v0.PATCH("/users/:id", userHandler.UpdateUser)
		v0.DELETE("/users/:id", userHandler.DeleteUser)
		// v0.GET("/users/:username", userHandler.GetUserByName)
		// Auth routes
		v0.POST("/auth/login", authHandler.Login)
		v0.POST("/auth/refresh", authHandler.Refresh)
		v0.POST("/auth/logout", middleware.JWTAuthMiddleware(), authHandler.Logout)
	}

	v1 := app.Group("/api/v1")
	v1.Use(middleware.JWTAuthMiddleware())
	{
		v1.GET("/orders", orderHandler.ListOrders)
		v1.GET("/orders/:id", orderHandler.GetOrder)
		v1.POST("/orders", orderHandler.CreateOrder)
		v1.PATCH("/orders/:id/status", orderHandler.UpdateOrderStatus)
		v1.DELETE("/orders/:id", orderHandler.DeleteOrder)
	}

	v2 := app.Group("/api/v2")
	v2.Use(middleware.JWTAuthMiddleware())
	{
		// Search routes
		v2.GET("/search/products", searchHandler.SearchProducts)
		v2.GET("/search/orders", searchHandler.SearchOrders)

		// List & create shops
		v2.GET("/shops", shopHandler.ListShops)
		v2.POST("/shops", shopHandler.CreateShop)

		// Product routes (more specific) MUST be registered before the
		// generic '/shops/:id' route to avoid gin wildcard conflicts.
		v2.POST("/shops/:id/products", shopHandler.CreateProduct)
		v2.GET("/shops/:id/products", shopHandler.ListProducts)
		v2.GET("/shops/:id/products/search", shopHandler.GetProductByName)

		// Product operations by product id
		v2.GET("/products/:id", shopHandler.GetProductByCode)
		v2.PATCH("/products/:id", shopHandler.UpdateProduct)
		v2.DELETE("/products/:id", shopHandler.DeleteProduct)
		v2.DELETE("/products", shopHandler.BatchDeleteProducts)

		// Shop detail/update/delete (generic param) registered after product routes
		v2.GET("/shops/:id", shopHandler.GetShop)
		v2.PATCH("/shops/:id", shopHandler.UpdateShop)
		v2.DELETE("/shops/:id", shopHandler.DeleteShop)
		v2.DELETE("/shops", shopHandler.BatchDeleteShops)
	}

	v3 := app.Group("api/v3")
	{
		v3.GET("/shops/:id/comments", commentHandler.ListCommentsByShop)
		v3.POST("/shops/:id/comments", middleware.JWTAuthMiddleware(), commentHandler.CreateComment)
		v3.DELETE("/shops/:id/comments", commentHandler.DeleteComments)
	}

}

func (app *Application) run() error {
	err := app.Run(app.config.address)
	if err != nil {
		return err
	}
	return nil
}
