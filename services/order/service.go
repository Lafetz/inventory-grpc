package order

import "github.com/google/uuid"

type service struct {
	repo repository
}

func (srv *service) GetOrder(id uuid.UUID) (*order, error) {
	return srv.repo.GetOrder(id)
}
func (srv *service) AddOrder(product *order) error {
	return srv.repo.AddOrder(product)
}
