package v1

import (
	"github.com/asliddinberdiev/i_tv_task/internal/modules/user"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Module("delivery_http_v1", fx.Provide(NewV1Routes))

type V1Routes struct {
	users user.Handler
}

type V1RoutesParams struct {
	fx.In
	Users user.Handler
}

func NewV1Routes(params V1RoutesParams) *V1Routes {
	return &V1Routes{
		users: params.Users,
	}
}

func (h *V1Routes) SetupRoutes(api *gin.RouterGroup) {
	v1Group := api.Group("/v1")
	{

		users := v1Group.Group("/users")
		{
			users.GET("/:id", h.users.GetUserID)
		}
	}
}
