package user

type Service interface {
	GetByID() error
}

type service struct {
	r Repository
}

func NewService(repository Repository) Service {
	return &service{r: repository}
}

func (s *service) GetByID() error {
	return s.r.GetByID()
}
