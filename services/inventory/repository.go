package inventory

import "github.com/google/uuid"

type store struct {
}

func (s *store) AddProduct(title string, desc string) (*product, error) {
	return NewProduct(title, desc), nil
}
func (s *store) GetProduct(id uuid.UUID) (*product, error) {
	return &product{Id: id, Title: "hello", Description: "world hello"}, nil
}
func (s *store) DecProduct(id uuid.UUID, quanitty int32) error {
	return nil
}
func NewDb() *store {
	return &store{}
}
