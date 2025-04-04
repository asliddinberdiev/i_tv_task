package v1

import (
	"github.com/asliddinberdiev/i_tv_task/internal/config"
	"github.com/asliddinberdiev/i_tv_task/internal/modules"
	logger "github.com/asliddinberdiev/i_tv_task/pkgs/logger/zap"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	log logger.Logger
	cfg *config.Config
	m   *modules.Modules
}

func NewHandler(log logger.Logger, cfg *config.Config, modules *modules.Modules) *Handler {
	return &Handler{log: log, cfg: cfg, m: modules}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initUserRoutes(v1)
	}
}
