package repository

import (
	"context"

	"github.com/steven230500/cartagena-api/internal/domain"
)

type RouteProgressRepository interface {
	Get(ctx context.Context, userID, routeID string) (domain.RouteProgress, error)
	Upsert(ctx context.Context, userID, routeID string, currentStep int) (domain.RouteProgress, error)
	CountCompleted(ctx context.Context, userID string) (int, error)
}
