package repository

import (
	"context"

	"github.com/steven230500/cartagena-api/internal/domain"
)

type BusinessClaimRepository interface {
	Create(ctx context.Context, businessID, userID, message string) (domain.BusinessClaim, error)
	List(ctx context.Context, filter domain.BusinessClaimFilter) ([]domain.BusinessClaim, error)
	Get(ctx context.Context, id string) (domain.BusinessClaim, error)
	Resolve(ctx context.Context, id, status string) (domain.BusinessClaim, error)
}
