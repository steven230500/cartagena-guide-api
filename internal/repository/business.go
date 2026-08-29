package repository

import (
	"context"

	"github.com/steven230500/cartagena-api/internal/domain"
)

type BusinessRepository interface {
	List(ctx context.Context, filter domain.BusinessFilter) ([]domain.Business, error)
	GetBySlug(ctx context.Context, slug string) (domain.Business, error)
	GetByID(ctx context.Context, id string) (domain.Business, error)
	ListByOwner(ctx context.Context, ownerID string) ([]domain.Business, error)
	Create(ctx context.Context, b domain.Business) (domain.Business, error)
	Update(ctx context.Context, id string, b domain.Business) (domain.Business, error)
	SetOwner(ctx context.Context, id, ownerID string) (domain.Business, error)
	Delete(ctx context.Context, id string) error
}
