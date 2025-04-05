package user

import (
	"github.com/asliddinberdiev/i_tv_task/internal/storage/postgres"
	"gorm.io/gorm"
)

type Repository interface {
	GetByID(id string) (UserID, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(psql postgres.PostgresDB) Repository {
	return &repository{db: psql.DB()}
}

func (r *repository) GetByID(id string) (UserID, error) {
	return UserID{ID: id}, nil
}
