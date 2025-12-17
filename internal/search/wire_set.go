package search

import (
	"github.com/google/wire"
	ordersearch "github.com/myproject/shop/internal/search/order"
	"github.com/myproject/shop/internal/search/product"
)

var ProviderSet = wire.NewSet(
	NewHandler,
	product.NewService,
	ordersearch.NewService,
)
