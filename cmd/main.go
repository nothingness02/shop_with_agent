package main

import (
	"log"

	config "github.com/myproject/shop/internal/config"
	"github.com/myproject/shop/pkg/middleware"
)

func main() {
	// 1. 加载配置
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// 2. 初始化 JWT Secret (从配置中读取)
	middleware.InitJWT(cfg.JWT.Secret)

	// 3. 依赖注入初始化应用 (调用 wire_gen.go 中的函数)
	app, err := InitializeApp(cfg)
	if err != nil {
		log.Fatalf("cannot initialize app: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
