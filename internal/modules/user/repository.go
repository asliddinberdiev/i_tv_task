package user

import (
	"github.com/asliddinberdiev/i_tv_task/internal/modules/common"
	"github.com/asliddinberdiev/i_tv_task/internal/storage/postgres"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user User) (*common.ResponseID, error)
	GetByEmail(email string) (*User, error)
	GetByID(req common.RequestID) (*User, error)
	Update(user User) (*common.ResponseID, error)
	Delete(req common.RequestID) (*common.ResponseID, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(psql postgres.PostgresDB) Repository {
	return &repository{db: psql.DB()}
}

func (r *repository) Create(user User) (*common.ResponseID, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &common.ResponseID{ID: user.ID}, nil
}

func (r *repository) GetByID(req common.RequestID) (*User, error) {
	var user User
	if err := r.db.Where("id = ?", req.ID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetByEmail(email string) (*User, error) {
	var user User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) Update(user User) (*common.ResponseID, error) {
	if err := r.db.Model(&user).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		return nil, err
	}
	return &common.ResponseID{ID: user.ID}, nil
}

func (r *repository) Delete(req common.RequestID) (*common.ResponseID, error) {
	if err := r.db.Delete(&User{}, req.ID).Error; err != nil {
		return nil, err
	}
	return &common.ResponseID{ID: req.ID}, nil
}
