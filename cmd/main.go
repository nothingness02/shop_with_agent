package main

import (
	"log"

	"github.com/myproject/shop/cmd/validator"
	config "github.com/myproject/shop/internal/config"
	"github.com/myproject/shop/pkg/logger"
	"github.com/myproject/shop/pkg/middleware"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	if err := logger.Init("logs/app.jsonl"); err != nil {
		log.Fatalf("cannot init logger: %v", err)
	}
	defer func() {
		_ = logger.Close()
	}()

	middleware.InitJWT(cfg.JWT.Secret)
	validator.RegisterPhoneValidator()
	app, err := InitializeApp(cfg)
	if err != nil {
		log.Fatalf("cannot initialize app: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}