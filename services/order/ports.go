package order

import (
	"context"

	"github.com/google/uuid"
)

type repository interface {
	AddOrder(*order) error
	GetOrder(uuid.UUID) (*order, error)
}
type serviceApi interface {
	AddOrder(context.Context, *order) error
	GetOrder(context.Context, uuid.UUID) (order, error)
}
