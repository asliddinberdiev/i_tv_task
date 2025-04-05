package user

import "go.uber.org/fx"

var Module = fx.Module(
	"user_module",
	fx.Provide(
		NewRepository,
		NewService,
		NewHandler,
	),
)
