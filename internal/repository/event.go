package repository

import (
	"context"

	"github.com/steven230500/cartagena-api/internal/domain"
)

type EventRepository interface {
	List(ctx context.Context) ([]domain.Event, error)
	GetBySlug(ctx context.Context, slug string) (domain.Event, error)
	Create(ctx context.Context, e domain.Event) (domain.Event, error)
	Update(ctx context.Context, id string, e domain.Event) (domain.Event, error)
	Delete(ctx context.Context, id string) error
}
