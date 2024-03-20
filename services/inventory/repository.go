package inventory

import "github.com/google/uuid"

type store struct {
}

func (s *store) AddProduct(p *product) error {
	return nil
}
func (s *store) GetProduct(id uuid.UUID) (*product, error) {
	return &product{Id: id, Title: "hello", Description: "world hello"}, nil
}
