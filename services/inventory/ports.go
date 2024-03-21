package inventory

import (
	"github.com/google/uuid"
)

type repository interface {
	AddProduct(string, string) (*product, error)
	GetProduct(uuid.UUID) (*product, error)
	DecProduct(uuid.UUID, int32) error
}
type serviceApi interface {
	AddProduct(string, string) (*product, error)
	GetProduct(uuid.UUID) (*product, error)
}
