package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initUserRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.GET("", h.test)
	}
}

func (h *Handler) test(c *gin.Context) {
	h.m.User.GetByID()
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
