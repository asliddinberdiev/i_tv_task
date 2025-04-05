package movie

import (
	"net/http"
	"strconv"

	"github.com/asliddinberdiev/i_tv_task/internal/modules/common"
	"github.com/asliddinberdiev/i_tv_task/pkgs/helper"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Handler interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	GetAll(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{service: service}
}

// @Summary Create a new movie
// @Description Create a new movie
// @Tags movies
// @Accept json
// @Produce json
// @Param movie body MovieCreateInput true "Movie"
// @Success 201 {object} common.ResponseID
// @Failure 400 {object} common.ResponseError
// @Failure 500 {object} common.ResponseError
// @Router /api/v1/movies [post]
func (h *handler) Create(c *gin.Context) {
	var req MovieCreateInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			common.ResponseError{
				Status:  http.StatusBadRequest,
				Message: "Invalid request",
			},
		)
		return
	}

	if err := common.Validate.Struct(req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			common.ResponseError{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			},
		)
		return
	}

	movie, err := h.service.Create(req)
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

		c.JSON(
			http.StatusInternalServerError,
			common.ResponseError{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusCreated,
		common.ResponseID{
			Status:  http.StatusCreated,
			Message: "Movie created successfully",
			ID:      movie.ID,
		},
	)
}

// @Summary Get a movie by ID
// @Description Get a movie by ID
// @Tags movies
// @Accept json
// @Produce json
// @Param id path string true "Movie ID"
// @Success 200 {object} common.Response
// @Failure 400 {object} common.ResponseError
// @Failure 404 {object} common.ResponseError
// @Failure 500 {object} common.ResponseError
// @Router /api/v1/movies/{id} [get]
func (h *handler) GetByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(
			http.StatusBadRequest,
			common.ResponseError{
				Status:  http.StatusBadRequest,
				Message: "Invalid movie id",
			},
		)
		return
	}

	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			common.ResponseError{
				Status:  http.StatusBadRequest,
				Message: "Invalid movie id",
			},
		)
		return
	}

	movie, err := h.service.GetByID(common.RequestID{ID: uint(idUint)})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(
				http.StatusNotFound,
				common.ResponseError{
					Status:  http.StatusNotFound,
					Message: "Movie not found",
				},
			)
			return
		}

		c.JSON(
			http.StatusInternalServerError,
			common.ResponseError{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		common.Response{
			Status:  http.StatusOK,
			Message: "Movie fetched successfully",
			Data:    movie,
		},
	)
}

// @Summary Get all movies
// @Description Get all movies
// @Tags movies
// @Accept json
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(10)
// @Param search query string false "Search"
// @Success 200 {object} common.ResponseWithList
// @Failure 500 {object} common.ResponseError
// @Router /api/v1/movies [get]
func (h *handler) GetAll(c *gin.Context) {
	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil || limit < 1 {
		limit = 10
	}

	var req common.RequestSearch
	req.Page = page
	req.Limit = limit
	req.Search = c.Query("search")

	movies, err := h.service.GetAll(req)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			common.ResponseError{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		common.ResponseWithList{
			Status:  http.StatusOK,
			Message: "Movies fetched successfully",
			Total:   movies.Total,
			Data:    movies.Movies,
		},
	)
}

// @Summary Update a movie by ID
// @Description Update a movie by ID
// @Tags movies
// @Accept json
// @Produce json
// @Param id path string true "Movie ID"
// @Param movie body MovieUpdateInput true "Movie"
// @Success 200 {object} common.ResponseID
// @Failure 400 {object} common.ResponseError
// @Failure 404 {object} common.ResponseError
// @Failure 500 {object} common.ResponseError
// @Router /api/v1/movies/{id} [put]
func (h *handler) Update(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(
			http.StatusBadRequest,
			common.ResponseError{
				Status:  http.StatusBadRequest,
				Message: "Invalid movie id",
			},
		)
		return
	}

	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			common.ResponseError{
				Status:  http.StatusBadRequest,
				Message: "Invalid movie id",
			},
		)
		return
	}

	var req MovieUpdateInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			common.ResponseError{
				Status:  http.StatusBadRequest,
				Message: "Invalid body",
			},
		)
		return
	}

	if err := common.Validate.Struct(req); err != nil {
		c.JSON(
			http.StatusBadRequest,
			common.ResponseError{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			},
		)
		return
	}

	req.ID = uint(idUint)
	movie, err := h.service.Update(req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(
				http.StatusNotFound,
				common.ResponseError{
					Status:  http.StatusNotFound,
					Message: "Movie not found",
				},
			)
			return
		}

		c.JSON(
			http.StatusInternalServerError,
			common.ResponseError{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		common.ResponseID{
			Status:  http.StatusOK,
			Message: "Movie updated successfully",
			ID:      movie.ID,
		},
	)
}

// @Summary Delete a movie by ID
// @Description Delete a movie by ID
// @Tags movies
// @Accept json
// @Produce json
// @Param id path string true "Movie ID"
// @Success 200 {object} common.ResponseID
// @Failure 400 {object} common.ResponseError
// @Failure 404 {object} common.ResponseError
// @Failure 500 {object} common.ResponseError
// @Router /api/v1/movies/{id} [delete]
func (h *handler) Delete(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(
			http.StatusBadRequest,
			common.ResponseError{
				Status:  http.StatusBadRequest,
				Message: "Invalid movie id",
			},
		)
		return
	}

	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			common.ResponseError{
				Status:  http.StatusBadRequest,
				Message: "Invalid movie id",
			},
		)
		return
	}

	movie, err := h.service.Delete(common.RequestID{ID: uint(idUint)})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(
				http.StatusNotFound,
				common.ResponseError{
					Status:  http.StatusNotFound,
					Message: "Movie not found",
				},
			)
			return
		}

		c.JSON(
			http.StatusInternalServerError,
			common.ResponseError{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		common.ResponseID{
			Status:  http.StatusOK,
			Message: "Movie deleted successfully",
			ID:      movie.ID,
		},
	)
}
