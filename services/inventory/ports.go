package inventory

import (
	"context"

	"github.com/google/uuid"
)

type repository interface {
	AddProduct(*Product) error
	GetProduct(uuid.UUID) (*Product, error)
}
type serviceApi interface {
	AddProduct(context.Context, *Product) error
	GetProduct(context.Context, uuid.UUID) (Product, error)
}
