package http

import (
	"net/http"

	"github.com/asliddinberdiev/i_tv_task/docs"
	v1 "github.com/asliddinberdiev/i_tv_task/internal/api/http/v1"
	"github.com/asliddinberdiev/i_tv_task/internal/config"
	"github.com/asliddinberdiev/i_tv_task/internal/modules"
	logger "github.com/asliddinberdiev/i_tv_task/pkgs/logger/zap"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	log     logger.Logger
	cfg     *config.Config
	modules *modules.Modules
}

func NewHandler(log logger.Logger, cfg *config.Config, modules *modules.Modules) *Handler {
	return &Handler{log: log, cfg: cfg, modules: modules}
}

// @title I_TV_TASK API
// @version v1
// @description REST API for I_TV_TASK APP

// @host localhost:8000
// @BasePath /api/v1/
func (h *Handler) Init() *gin.Engine {
	if h.cfg.App.Environment != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.Use(corsMiddleware)

	docs.SwaggerInfo.Host = h.cfg.GetAppAddr()
	if h.cfg.App.Environment == "dev" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.log, h.cfg, h.modules)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
