package v1

import (
	"github.com/asliddinberdiev/i_tv_task/internal/modules/movie"
	"github.com/asliddinberdiev/i_tv_task/internal/modules/user"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Module("delivery_http_v1", fx.Provide(NewV1Routes))

type V1Routes struct {
	users  user.Handler
	movies movie.Handler
}

type V1RoutesParams struct {
	fx.In
	Users  user.Handler
	Movies movie.Handler
}

func NewV1Routes(params V1RoutesParams) *V1Routes {
	return &V1Routes{
		users:  params.Users,
		movies: params.Movies,
	}
}

func (h *V1Routes) SetupRoutes(api *gin.RouterGroup) {
	v1Group := api.Group("/v1")
	{
		users := v1Group.Group("/users")
		{
			users.POST("/register", h.users.Register)
			users.POST("/login", h.users.Login)
		}

		movies := v1Group.Group("/movies")
		{
			movies.POST("", h.movies.Create)
			movies.GET("/:id", h.movies.GetByID)
			movies.GET("", h.movies.GetAll)
			movies.PUT("/:id", h.movies.Update)
			movies.DELETE("/:id", h.movies.Delete)
		}
	}
}
