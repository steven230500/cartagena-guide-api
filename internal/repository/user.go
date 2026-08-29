package repository

import (
	"context"

	"github.com/steven230500/cartagena-api/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, email, passwordHash string) (domain.User, error)
	// GetByEmail devuelve también el hash — es lo único que Login necesita para comparar.
	GetByEmail(ctx context.Context, email string) (domain.User, string, error)
	GetByID(ctx context.Context, id string) (domain.User, error)
	UpdateRole(ctx context.Context, id, role string) (domain.User, error)
}
