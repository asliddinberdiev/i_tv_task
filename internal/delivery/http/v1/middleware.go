package v1

import (
	"net/http"

	"github.com/asliddinberdiev/i_tv_task/internal/config"
	"github.com/asliddinberdiev/i_tv_task/internal/modules/common"
	"github.com/asliddinberdiev/i_tv_task/internal/modules/user"
	"github.com/asliddinberdiev/i_tv_task/pkgs/auth"
	"github.com/asliddinberdiev/i_tv_task/pkgs/helper"
	"github.com/gin-gonic/gin"
)

func jwtMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				common.ResponseError{
					Status:  http.StatusUnauthorized,
					Message: "header is required",
				},
			)
			return
		}

		var claims user.UserClaims
		if err := auth.ParseToken(token, cfg.Auth.SecretKey, &claims); err != nil {
			if helper.ErrorIs(err, "token is expired") {
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					common.ResponseError{
						Status:  http.StatusUnauthorized,
						Message: err.Error(),
					},
				)
				return
			}
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				common.ResponseError{
					Status:  http.StatusBadRequest,
					Message: "Invalid token",
				},
			)
			return
		}

		c.Set("user_id", claims.ID)
		c.Next()
	}
}
