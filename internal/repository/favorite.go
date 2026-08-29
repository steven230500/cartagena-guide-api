package repository

import (
	"context"

	"github.com/steven230500/cartagena-api/internal/domain"
)

type FavoriteRepository interface {
	List(ctx context.Context, userID string) ([]domain.Business, error)
	Add(ctx context.Context, userID, businessID string) error
	Remove(ctx context.Context, userID, businessID string) error
}
