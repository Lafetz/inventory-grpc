package inventory

import (
	"context"

	"github.com/google/uuid"
)

type repository interface {
	AddProduct(*product) error
	GetProduct(uuid.UUID) (*product, error)
}
type serviceApi interface {
	AddProduct(context.Context, *product) error
	GetProduct(context.Context, uuid.UUID) (product, error)
}
