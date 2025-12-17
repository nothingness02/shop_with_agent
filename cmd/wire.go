package main

import (
	"log"

	"github.com/google/wire"
	auth "github.com/myproject/shop/internal/Auth"
	comment "github.com/myproject/shop/internal/Comment"
	"github.com/myproject/shop/internal/Order"
	"gorm.io/gorm"

	shop "github.com/myproject/shop/internal/Shop"
	user "github.com/myproject/shop/internal/User"
	config "github.com/myproject/shop/internal/config"
	"github.com/myproject/shop/internal/search"
	"github.com/myproject/shop/pkg/database"
	"github.com/myproject/shop/pkg/middleware"
)

func provideRedisStore(cfg *config.Config) *middleware.RedisStore {
	store := middleware.NewRedisStore(
		cfg.Redis.Addr,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	middleware.InitRedis(store.Client)
	return store
}

func provideDB(cfg *config.Config) (*database.Database, error) {
	db, err := database.NewDB(cfg.Database.BuildPostgresDSN(""))
	if err != nil {
		log.Fatal(err)
		return db, err
	}
	if err := db.AutoMigrate(&Order.Order{}, &Order.OrderItem{},
		&shop.Shop{}, &shop.Product{},
		&shop.Category{}, &user.User{}, &comment.Comment{},
	); err != nil {
		log.Fatal(err)
		return db, err
	}
	return db, err

}

func provideGormDB(db *database.Database) *gorm.DB {
	return db.DB
}

// InitializeApp 是我们要生成的“总构造函数”
func InitializeApp(cfg *config.Config) (*Application, error) {
	wire.Build(
		// 1. 基础设施
		provideDB,
		provideRedisStore,
		provideGormDB,
		user.ProviderSet,
		auth.ProviderSet,
		Order.ProviderSet,
		shop.ProviderSet,
		search.ProviderSet,
		comment.ProviderSet,
		NewApplication,
	)
	return &Application{}, nil
}
