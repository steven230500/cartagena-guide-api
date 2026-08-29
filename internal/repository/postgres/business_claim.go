package postgres

import (
	"context"
	"time"

	"github.com/steven230500/cartagena-api/internal/db"
	"github.com/steven230500/cartagena-api/internal/domain"
)

type BusinessClaimRepository struct {
	Q *db.Queries
}

func NewBusinessClaimRepository(q *db.Queries) *BusinessClaimRepository {
	return &BusinessClaimRepository{Q: q}
}

func strOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func toDomainClaim(c db.BusinessClaim) domain.BusinessClaim {
	var resolvedAt *time.Time
	if c.ResolvedAt.Valid {
		t := c.ResolvedAt.Time
		resolvedAt = &t
	}
	return domain.BusinessClaim{
		ID: c.ID.String(), BusinessID: c.BusinessID.String(), UserID: c.UserID.String(),
		Message: strOrEmpty(c.Message), Status: c.Status,
		CreatedAt: c.CreatedAt.Time, ResolvedAt: resolvedAt,
	}
}

func toDomainClaimRow(c db.ListBusinessClaimsRow) domain.BusinessClaim {
	var resolvedAt *time.Time
	if c.ResolvedAt.Valid {
		t := c.ResolvedAt.Time
		resolvedAt = &t
	}
	return domain.BusinessClaim{
		ID: c.ID.String(), BusinessID: c.BusinessID.String(), BusinessName: c.BusinessName, BusinessSlug: c.BusinessSlug,
		UserID: c.UserID.String(), UserEmail: c.UserEmail, Message: strOrEmpty(c.Message), Status: c.Status,
		CreatedAt: c.CreatedAt.Time, ResolvedAt: resolvedAt,
	}
}

func (r *BusinessClaimRepository) Create(ctx context.Context, businessID, userID, message string) (domain.BusinessClaim, error) {
	bid, err := mustUUID(businessID)
	if err != nil {
		return domain.BusinessClaim{}, err
	}
	uid, err := mustUUID(userID)
	if err != nil {
		return domain.BusinessClaim{}, err
	}
	var msg *string
	if message != "" {
		msg = &message
	}
	c, err := r.Q.CreateBusinessClaim(ctx, db.CreateBusinessClaimParams{BusinessID: bid, UserID: uid, Message: msg})
	if err != nil {
		return domain.BusinessClaim{}, err
	}
	return toDomainClaim(c), nil
}

func (r *BusinessClaimRepository) Get(ctx context.Context, id string) (domain.BusinessClaim, error) {
	cid, err := mustUUID(id)
	if err != nil {
		return domain.BusinessClaim{}, err
	}
	c, err := r.Q.GetBusinessClaim(ctx, cid)
	if err != nil {
		return domain.BusinessClaim{}, domain.ErrNotFound
	}
	return toDomainClaim(c), nil
}

func (r *BusinessClaimRepository) Resolve(ctx context.Context, id, status string) (domain.BusinessClaim, error) {
	cid, err := mustUUID(id)
	if err != nil {
		return domain.BusinessClaim{}, err
	}
	c, err := r.Q.ResolveBusinessClaim(ctx, db.ResolveBusinessClaimParams{ID: cid, Status: status})
	if err != nil {
		return domain.BusinessClaim{}, domain.ErrNotFound
	}
	return toDomainClaim(c), nil
}

func (r *BusinessClaimRepository) List(ctx context.Context, filter domain.BusinessClaimFilter) ([]domain.BusinessClaim, error) {
	params := db.ListBusinessClaimsParams{Status: strOrNil(filter.Status)}
	if filter.UserID != "" {
		uid, err := mustUUID(filter.UserID)
		if err != nil {
			return nil, err
		}
		params.UserID = uid
	}
	// filter.UserID == "" deja params.UserID en su zero value (Valid:false),
	// que sqlc.narg('user_id') ya interpreta como "sin filtro".
	rows, err := r.Q.ListBusinessClaims(ctx, params)
	if err != nil {
		return nil, err
	}
	out := make([]domain.BusinessClaim, len(rows))
	for i, c := range rows {
		out[i] = toDomainClaimRow(c)
	}
	return out, nil
}
