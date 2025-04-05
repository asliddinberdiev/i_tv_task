package movie

import (
	"github.com/asliddinberdiev/i_tv_task/internal/modules/common"
	"gorm.io/gorm"
)

type Service interface {
	Create(req MovieCreateInput) (*common.ResponseID, error)
	GetByID(req common.RequestID) (*MovieResponse, error)
	GetAll(req common.RequestSearch) (*MovieListResponse, error)
	Update(req MovieUpdateInput) (*common.ResponseID, error)
	Delete(req common.RequestID) (*common.ResponseID, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(req MovieCreateInput) (*common.ResponseID, error) {
	return s.repo.Create(Movie{
		Title:    req.Title,
		Year:     req.Year,
		Genre:    req.Genre,
		Rating:   req.Rating,
		Director: req.Director,
	})
}

func (s *service) GetByID(req common.RequestID) (*MovieResponse, error) {
	movie, err := s.repo.GetByID(req)
	if err != nil {
		return nil, err
	}
	return &MovieResponse{
		ID:        movie.ID,
		Title:     movie.Title,
		Year:      movie.Year,
		Genre:     movie.Genre,
		Rating:    movie.Rating,
		Director:  movie.Director,
		CreatedAt: movie.CreatedAt,
		UpdatedAt: movie.UpdatedAt,
	}, nil
}

func (s *service) GetAll(req common.RequestSearch) (*MovieListResponse, error) {
	return s.repo.GetAll(req)
}

func (s *service) Update(req MovieUpdateInput) (*common.ResponseID, error) {
	return s.repo.Update(Movie{
		Model:    gorm.Model{ID: req.ID},
		Title:    req.Title,
		Year:     req.Year,
		Genre:    req.Genre,
		Rating:   req.Rating,
		Director: req.Director,
	})
}

func (s *service) Delete(req common.RequestID) (*common.ResponseID, error) {
	return s.repo.Delete(req)
}
