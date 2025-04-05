package user

import (
	"github.com/asliddinberdiev/i_tv_task/internal/modules/common"
)

type Service interface {
	Create(req User) (*common.ResponseID, error)
	GetByEmail(email string) (*User, error)
	GetByID(req common.RequestID) (*User, error)
	Update(user User) (*common.ResponseID, error)
	Delete(req common.RequestID) (*common.ResponseID, error)
}

type service struct {
	r Repository
}

func NewService(repository Repository) Service {
	return &service{r: repository}
}

func (s *service) Create(req User) (*common.ResponseID, error) {
	return s.r.Create(req)
}

func (s *service) GetByEmail(email string) (*User, error) {
	return s.r.GetByEmail(email)
}

func (s *service) GetByID(req common.RequestID) (*User, error) {
	return s.r.GetByID(req)
}

func (s *service) Update(user User) (*common.ResponseID, error) {
	return s.r.Update(user)
}

func (s *service) Delete(req common.RequestID) (*common.ResponseID, error) {
	return s.r.Delete(req)
}
