package postgres

import (
	"context"

	"github.com/steven230500/cartagena-api/internal/db"
	"github.com/steven230500/cartagena-api/internal/domain"
)

type FavoriteRepository struct {
	Q *db.Queries
}

func NewFavoriteRepository(q *db.Queries) *FavoriteRepository {
	return &FavoriteRepository{Q: q}
}

func (r *FavoriteRepository) List(ctx context.Context, userID string) ([]domain.Business, error) {
	uid, err := mustUUID(userID)
	if err != nil {
		return nil, err
	}
	rows, err := r.Q.ListFavoriteBusinesses(ctx, uid)
	if err != nil {
		return nil, err
	}
	out := make([]domain.Business, len(rows))
	for i, b := range rows {
		out[i] = toDomainBusiness(b)
	}
	return out, nil
}

func (r *FavoriteRepository) Add(ctx context.Context, userID, businessID string) error {
	uid, err := mustUUID(userID)
	if err != nil {
		return err
	}
	bid, err := mustUUID(businessID)
	if err != nil {
		return err
	}
	return r.Q.AddFavorite(ctx, db.AddFavoriteParams{UserID: uid, BusinessID: bid})
}

func (r *FavoriteRepository) Remove(ctx context.Context, userID, businessID string) error {
	uid, err := mustUUID(userID)
	if err != nil {
		return err
	}
	bid, err := mustUUID(businessID)
	if err != nil {
		return err
	}
	return r.Q.RemoveFavorite(ctx, db.RemoveFavoriteParams{UserID: uid, BusinessID: bid})
}
