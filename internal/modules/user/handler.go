package user

import (
	"net/http"

	"github.com/asliddinberdiev/i_tv_task/internal/config"
	logger "github.com/asliddinberdiev/i_tv_task/pkgs/logger/zap"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	GetUserID(c *gin.Context)
}

type handler struct {
	s   Service
	log logger.Logger
	cfg *config.Config
}

func NewHandler(service Service, log logger.Logger, cfg *config.Config) Handler {
	return &handler{s: service, log: log, cfg: cfg}
}

// GetUserID godoc
// @Summary Get user id
// @Description Get user id
// @Accept json
// @Produce json
// @Success 200 {object} user.UserID
// @Router /user/{id} [get]
func (h *handler) GetUserID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		h.log.Error("GetUser", logger.String("error", "id is required"))
		c.JSON(http.StatusBadRequest, gin.H{"message": "id is required"})
		return
	}

	response, err := h.s.GetByID(id)
	if err != nil {
		h.log.Error("GetUser", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, response)
}
