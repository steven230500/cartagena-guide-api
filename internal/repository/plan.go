package repository

import (
	"context"

	"github.com/steven230500/cartagena-api/internal/domain"
)

type PlanRepository interface {
	List(ctx context.Context) ([]domain.Plan, error)
	Create(ctx context.Context, p domain.Plan) (domain.Plan, error)
	Update(ctx context.Context, id string, p domain.Plan) (domain.Plan, error)
	Delete(ctx context.Context, id string) error
}
