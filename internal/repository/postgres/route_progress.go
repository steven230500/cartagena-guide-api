package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/steven230500/cartagena-api/internal/db"
	"github.com/steven230500/cartagena-api/internal/domain"
)

type RouteProgressRepository struct {
	Q *db.Queries
}

func NewRouteProgressRepository(q *db.Queries) *RouteProgressRepository {
	return &RouteProgressRepository{Q: q}
}

func (r *RouteProgressRepository) Get(ctx context.Context, userID, routeID string) (domain.RouteProgress, error) {
	uid, err := mustUUID(userID)
	if err != nil {
		return domain.RouteProgress{}, err
	}
	rid, err := mustUUID(routeID)
	if err != nil {
		return domain.RouteProgress{}, err
	}
	row, err := r.Q.GetRouteProgress(ctx, db.GetRouteProgressParams{UserID: uid, RouteID: rid})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.RouteProgress{}, domain.ErrNotFound
		}
		return domain.RouteProgress{}, err
	}
	return toDomainRouteProgress(row), nil
}

func (r *RouteProgressRepository) Upsert(ctx context.Context, userID, routeID string, currentStep int) (domain.RouteProgress, error) {
	uid, err := mustUUID(userID)
	if err != nil {
		return domain.RouteProgress{}, err
	}
	rid, err := mustUUID(routeID)
	if err != nil {
		return domain.RouteProgress{}, err
	}
	row, err := r.Q.UpsertRouteProgress(ctx, db.UpsertRouteProgressParams{
		UserID:      uid,
		RouteID:     rid,
		CurrentStep: int32(currentStep),
	})
	if err != nil {
		return domain.RouteProgress{}, err
	}
	return toDomainRouteProgress(row), nil
}

func (r *RouteProgressRepository) CountCompleted(ctx context.Context, userID string) (int, error) {
	uid, err := mustUUID(userID)
	if err != nil {
		return 0, err
	}
	count, err := r.Q.CountCompletedRoutes(ctx, uid)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func toDomainRouteProgress(row db.UserRouteProgress) domain.RouteProgress {
	updatedAt := ""
	if row.UpdatedAt.Valid {
		updatedAt = row.UpdatedAt.Time.Format(time.RFC3339)
	}
	routeID := ""
	if p := uuidToPtr(row.RouteID); p != nil {
		routeID = *p
	}
	return domain.RouteProgress{
		RouteID:     routeID,
		CurrentStep: int(row.CurrentStep),
		UpdatedAt:   updatedAt,
	}
}
