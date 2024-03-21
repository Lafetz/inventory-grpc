package order

import "github.com/google/uuid"

type order struct {
	Id        uuid.UUID
	ProductId uuid.UUID
	Price     float32
	Quantity  int32
	Status    string
}
