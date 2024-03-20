package inventory

import "github.com/google/uuid"

type service struct {
	repo repository
}

func (srv *service) GetProduct(id uuid.UUID) (*Product, error) {
	return srv.repo.GetProduct(id)
}
func (srv *service) AddProduct(product *Product) error {
	return srv.repo.AddProduct(product)
}
