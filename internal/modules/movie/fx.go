package movie

import "go.uber.org/fx"

var Module = fx.Module(
	"movie_module",
	fx.Provide(
		NewRepository,
		NewService,
		NewHandler,
	),
)
