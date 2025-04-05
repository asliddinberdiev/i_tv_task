package user

import (
	"net/http"
	"time"

	"github.com/asliddinberdiev/i_tv_task/internal/config"
	"github.com/asliddinberdiev/i_tv_task/internal/modules/common"
	"github.com/asliddinberdiev/i_tv_task/pkgs/auth"
	"github.com/asliddinberdiev/i_tv_task/pkgs/helper"
	logger "github.com/asliddinberdiev/i_tv_task/pkgs/logger/zap"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Handler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type handler struct {
	s   Service
	log logger.Logger
	cfg *config.Config
}

func NewHandler(service Service, log logger.Logger, cfg *config.Config) Handler {
	return &handler{s: service, log: log, cfg: cfg}
}


// @Summary Register
// @Description Register
// @Tags auth
// @Accept json
// @Produce json
// @Param user body user.RegisterInput true "User"
// @Success 201 {object} user.TokenResponse
// @Failure 400 {object} common.ResponseError
// @Failure 500 {object} common.ResponseError
// @Router /api/v1/users/register [post]
func (h *handler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.log.Error("Register", logger.Error(err))
		c.JSON(
			http.StatusBadRequest,
			common.ResponseError{
				Status:  http.StatusBadRequest,
				Message: "Invalid body",
			},
		)
		return
	}

	if err := common.Validate.Struct(input); err != nil {
		c.JSON(
			http.StatusBadRequest,
			common.ResponseError{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			},
		)
		return
	}

	hashPassword, err := helper.PasswordHash(input.Password)
	if err != nil {
		h.log.Error("Register", logger.Error(err))
		c.JSON(
			http.StatusInternalServerError,
			common.ResponseError{
				Status:  http.StatusInternalServerError,
				Message: "Failed to hash password",
			},
		)
		return
	}
	newUser := User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  hashPassword,
	}

	user, err := h.s.Create(newUser)
	if err != nil {
		if helper.ErrorIs(err, "duplicate") {
			c.JSON(
				http.StatusBadRequest,
				common.ResponseError{
					Status:  http.StatusBadRequest,
					Message: "Already exists",
				},
			)
			return
		}

		h.log.Error("Register", logger.Error(err))
		c.JSON(
			http.StatusInternalServerError,
			common.ResponseError{
				Status:  http.StatusInternalServerError,
				Message: "Failed to create user",
			},
		)
		return
	}

	accessClaims := UserClaims{
		ID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(h.cfg.Auth.AccessTTL) * time.Second).Unix(),
		},
	}

	accessToken, err := auth.GenerateToken(accessClaims, h.cfg.Auth.SecretKey)
	if err != nil {
		h.log.Error("Login", logger.Error(err))
		c.JSON(
			http.StatusInternalServerError,
			common.ResponseError{
				Status:  http.StatusInternalServerError,
				Message: "Failed to generate token",
			},
		)
		return
	}

	refreshClaims := UserClaims{
		ID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(h.cfg.Auth.RefreshTTL) * time.Second).Unix(),
		},
	}
	refreshToken, err := auth.GenerateToken(refreshClaims, h.cfg.Auth.SecretKey)
	if err != nil {
		h.log.Error("Login", logger.Error(err))
		c.JSON(
			http.StatusInternalServerError,
			common.ResponseError{
				Status:  http.StatusInternalServerError,
				Message: "Failed to generate token",
			},
		)
		return
	}

	res := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ResponseID: common.ResponseID{
			Status:  http.StatusCreated,
			Message: "User created successfully",
			ID:      user.ID,
		},
	}

	c.JSON(http.StatusCreated, res)
}

// @Summary Login
// @Description Login
// @Tags auth
// @Accept json
// @Produce json
// @Param user body user.LoginInput true "User"
// @Success 200 {object} user.TokenResponse
// @Failure 400 {object} common.ResponseError
// @Failure 500 {object} common.ResponseError
// @Router /api/v1/users/login [post]
func (h *handler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.log.Error("Login", logger.Error(err))
		c.JSON(
			http.StatusBadRequest,
			common.ResponseError{
				Status:  http.StatusBadRequest,
				Message: "Invalid body",
			},
		)
		return
	}

	if err := common.Validate.Struct(input); err != nil {
		c.JSON(
			http.StatusBadRequest,
			common.ResponseError{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			},
		)
		return
	}

	user, err := h.s.GetByEmail(input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.log.Error("Login", logger.Error(err))
			c.JSON(
				http.StatusBadRequest,
				common.ResponseError{
					Status:  http.StatusBadRequest,
					Message: "Wrong email or password",
				},
			)
			return
		}

		h.log.Error("Login", logger.Error(err))
		c.JSON(
			http.StatusInternalServerError,
			common.ResponseError{
				Status:  http.StatusInternalServerError,
				Message: "Failed to get user",
			},
		)
		return
	}

	if !helper.PasswordCompare(user.Password, input.Password) {
		h.log.Error("Login", logger.Error(err))
		c.JSON(
			http.StatusBadRequest,
			common.ResponseError{
				Status:  http.StatusBadRequest,
				Message: "Wrong email or password",
			},
		)
		return
	}

	accessClaims := UserClaims{
		ID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(h.cfg.Auth.AccessTTL) * time.Second).Unix(),
		},
	}

	accessToken, err := auth.GenerateToken(accessClaims, h.cfg.Auth.SecretKey)
	if err != nil {
		h.log.Error("Login", logger.Error(err))
		c.JSON(
			http.StatusInternalServerError,
			common.ResponseError{
				Status:  http.StatusInternalServerError,
				Message: "Failed to generate token",
			},
		)
		return
	}

	refreshClaims := UserClaims{
		ID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(h.cfg.Auth.RefreshTTL) * time.Second).Unix(),
		},
	}
	refreshToken, err := auth.GenerateToken(refreshClaims, h.cfg.Auth.SecretKey)
	if err != nil {
		h.log.Error("Login", logger.Error(err))
		c.JSON(
			http.StatusInternalServerError,
			common.ResponseError{
				Status:  http.StatusInternalServerError,
				Message: "Failed to generate token",
			},
		)
		return
	}

	res := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ResponseID: common.ResponseID{
			Status:  http.StatusOK,
			Message: "User logged in successfully",
			ID:      user.ID,
		},
	}

	c.JSON(http.StatusOK, res)
}
