package inventory

import "github.com/google/uuid"

type service struct {
	repo repository
}

func (srv *service) DecProduct(id uuid.UUID, quanitty int32) error {
	return srv.repo.DecProduct(id, quanitty)
}
func (srv *service) GetProduct(id uuid.UUID) (*product, error) {
	return srv.repo.GetProduct(id)
}
func (srv *service) AddProduct(title string, desc string) (*product, error) {
	return srv.repo.AddProduct(title, desc)
}
func NewService(repo repository) *service {
	return &service{
		repo: repo,
	}
}
