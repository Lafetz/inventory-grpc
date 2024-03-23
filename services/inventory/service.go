package inventory

import (
	"context"

	"github.com/google/uuid"
)

type service struct {
	repo repository
}

func (srv *service) GetProduct(ctx context.Context, id uuid.UUID) (*product, error) {
	return srv.repo.GetProduct(ctx, id)
}
func (srv *service) AddProduct(ctx context.Context, title string, desc string) (*product, error) {
	return srv.repo.AddProduct(ctx, title, desc)
}
func NewService(repo repository) *service {
	return &service{
		repo: repo,
	}
}
