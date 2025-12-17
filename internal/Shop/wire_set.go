package shop

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewRepository,
	NewShopService,
	NewShopHandler,
)
