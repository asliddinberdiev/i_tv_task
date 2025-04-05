package user

type Service interface {
	GetByID(id string) (UserID, error)
}

type service struct {
	r Repository
}

func NewService(repository Repository) Service {
	return &service{r: repository}
}

func (s *service) GetByID(id string) (UserID, error) {
	return s.r.GetByID(id)
}
