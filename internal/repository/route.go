package repository

import (
	"context"

	"github.com/steven230500/cartagena-api/internal/domain"
)

type RouteRepository interface {
	List(ctx context.Context) ([]domain.Route, error)
	GetBySlug(ctx context.Context, slug string) (domain.Route, error)
	Create(ctx context.Context, r domain.Route) (domain.Route, error)
	Update(ctx context.Context, id string, r domain.Route) (domain.Route, error)
	Delete(ctx context.Context, id string) error
}
