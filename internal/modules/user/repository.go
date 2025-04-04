package user

import "gorm.io/gorm"

type Repository interface {
	GetByID() error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetByID() error {
	return nil
}
