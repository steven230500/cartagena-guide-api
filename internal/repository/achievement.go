package repository

import (
	"context"

	"github.com/steven230500/cartagena-api/internal/domain"
)

type AchievementRepository interface {
	List(ctx context.Context) ([]domain.Achievement, error)
	Create(ctx context.Context, a domain.Achievement) (domain.Achievement, error)
	Update(ctx context.Context, id string, a domain.Achievement) (domain.Achievement, error)
	Delete(ctx context.Context, id string) error
}
