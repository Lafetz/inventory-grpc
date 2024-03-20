package order

import "github.com/google/uuid"

type store struct {
}

func (s *store) AddOrder(p *order) error {
	return nil
}
func (s *store) GetOrder(id uuid.UUID) (*order, error) {
	return &order{Id: id, ProductId: uuid.New(), Price: 90, Quantity: 10}, nil
}
