package comment

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewRepository,
	NewCommentService,
	NewCommentHandler,
)
