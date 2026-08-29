package repository

import (
	"context"

	"github.com/steven230500/cartagena-api/internal/domain"
)

type ParishRepository interface {
	List(ctx context.Context) ([]domain.Parish, error)
	GetByID(ctx context.Context, id string) (domain.Parish, error)
	Create(ctx context.Context, p domain.Parish) (domain.Parish, error)
	Update(ctx context.Context, id string, p domain.Parish) (domain.Parish, error)
	Delete(ctx context.Context, id string) error
}
