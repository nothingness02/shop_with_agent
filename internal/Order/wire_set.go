package Order

import "github.com/google/wire"

// ProviderSet 把 Repository, Service, Handler 打包暴露给外部
var ProviderSet = wire.NewSet(
	NewRepository,
	NewOrderService,
	NewOrderHandler,
)
