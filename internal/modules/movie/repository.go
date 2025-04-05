package movie

import (
	"fmt"

	"github.com/asliddinberdiev/i_tv_task/internal/modules/common"
	"github.com/asliddinberdiev/i_tv_task/internal/storage/postgres"
	"gorm.io/gorm"
)

type Repository interface {
	Create(req Movie) (*common.ResponseID, error)
	GetByID(req common.RequestID) (*Movie, error)
	GetAll(req common.RequestSearch) (*MovieListResponse, error)
	Update(req Movie) (*common.ResponseID, error)
	Delete(req common.RequestID) (*common.ResponseID, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(psql postgres.PostgresDB) Repository {
	return &repository{db: psql.DB()}
}

func (r *repository) Create(req Movie) (*common.ResponseID, error) {
	if err := r.db.Create(&req).Error; err != nil {
		return nil, err
	}
	return &common.ResponseID{ID: req.ID}, nil
}

func (r *repository) GetByID(req common.RequestID) (*Movie, error) {
	var movie Movie
	if err := r.db.First(&movie, req.ID).Error; err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *repository) GetAll(req common.RequestSearch) (*MovieListResponse, error) {
	movies := make([]MovieResponse, 0)
	total := int64(0)

	tr := r.db.Begin()
	defer tr.Commit()

	if err := tr.Raw("SELECT COUNT(*) FROM movies").
		Scan(&total).Error; err != nil {
		tr.Rollback()
		return nil, err
	}

	query := `
		SELECT id, title, year, genre, rating, director, created_at, updated_at
		FROM movies
	`

	filter := ""
	if req.Search != "" {
		filter = fmt.Sprintf("WHERE title LIKE '%%%s%%'", req.Search)
	}

	query += filter
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT %d OFFSET %d", req.Limit, (req.Page-1)*req.Limit)

	if err := tr.Raw(query).Scan(&movies).Error; err != nil {
		tr.Rollback()
		return nil, err
	}

	response := MovieListResponse{
		Movies: movies,
		Total:  uint64(total),
	}

	return &response, nil
}

func (r *repository) Update(req Movie) (*common.ResponseID, error) {
	var movie Movie
	if err := r.db.Model(&movie).Where("id = ?", req.ID).Updates(req).Error; err != nil {
		return nil, err
	}
	return &common.ResponseID{ID: movie.ID}, nil
}

func (r *repository) Delete(req common.RequestID) (*common.ResponseID, error) {
	if err := r.db.Delete(&Movie{}, req.ID).Error; err != nil {
		return nil, err
	}
	return &common.ResponseID{ID: req.ID}, nil
}
