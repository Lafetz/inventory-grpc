package inventory

import (
	"context"

	"github.com/google/uuid"
)

type repository interface {
	AddProduct(context.Context, string, string) (*product, error)
	GetProduct(context.Context, uuid.UUID) (*product, error)
}
type serviceApi interface {
	AddProduct(context.Context, string, string) (*product, error)
	GetProduct(context.Context, uuid.UUID) (*product, error)
}
