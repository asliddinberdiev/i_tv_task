package http

import (
	"github.com/asliddinberdiev/i_tv_task/docs"
	"github.com/asliddinberdiev/i_tv_task/internal/config"
	v1 "github.com/asliddinberdiev/i_tv_task/internal/delivery/http/v1"
	logger "github.com/asliddinberdiev/i_tv_task/pkgs/logger/zap"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
)

var Module = fx.Module("delivery_http", fx.Provide(NewHandler))

type Handler struct {
	Router *gin.Engine
	cfg    *config.Config
	log    logger.Logger
}

type HandlerParams struct {
	fx.In

	Cfg *config.Config
	Log logger.Logger
	V1  *v1.V1Routes
}

func NewHandler(params HandlerParams) *Handler {
	router := gin.New()

	handler := &Handler{
		Router: router,
		cfg:    params.Cfg,
		log:    params.Log,
	}

	handler.Setup(params.V1)

	return handler
}

// @title I_TV API
// @version v1
// @description REST API for I_TV App

// @host localhost:8000
// @BasePath /api/v1/
func (h *Handler) Setup(v1Routes *v1.V1Routes) {
	h.Router.Use(
		gin.Recovery(),
		h.loggingMiddleware(),
		h.corsMiddleware(),
	)

	if h.cfg.App.Environment == "dev" {
		docs.SwaggerInfo.Host = h.cfg.GetAppAddr()
		h.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	api := h.Router.Group("/api")
	{
		v1Routes.SetupRoutes(api)
	}
}
