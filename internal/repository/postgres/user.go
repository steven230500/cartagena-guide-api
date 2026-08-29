package postgres

import (
	"context"

	"github.com/steven230500/cartagena-api/internal/db"
	"github.com/steven230500/cartagena-api/internal/domain"
)

type UserRepository struct {
	Q *db.Queries
}

func NewUserRepository(q *db.Queries) *UserRepository {
	return &UserRepository{Q: q}
}

func toDomainUser(u db.User) domain.User {
	return domain.User{ID: u.ID.String(), Email: u.Email, Role: u.Role, CreatedAt: u.CreatedAt.Time}
}

func (r *UserRepository) Create(ctx context.Context, email, passwordHash string) (domain.User, error) {
	u, err := r.Q.CreateUser(ctx, db.CreateUserParams{Email: email, PasswordHash: passwordHash})
	if err != nil {
		return domain.User{}, err
	}
	return toDomainUser(u), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, string, error) {
	u, err := r.Q.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, "", domain.ErrNotFound
	}
	return toDomainUser(u), u.PasswordHash, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (domain.User, error) {
	uid, err := mustUUID(id)
	if err != nil {
		return domain.User{}, err
	}
	u, err := r.Q.GetUserByID(ctx, uid)
	if err != nil {
		return domain.User{}, domain.ErrNotFound
	}
	return toDomainUser(u), nil
}

func (r *UserRepository) UpdateRole(ctx context.Context, id, role string) (domain.User, error) {
	uid, err := mustUUID(id)
	if err != nil {
		return domain.User{}, err
	}
	u, err := r.Q.UpdateUserRole(ctx, db.UpdateUserRoleParams{ID: uid, Role: role})
	if err != nil {
		return domain.User{}, domain.ErrNotFound
	}
	return toDomainUser(u), nil
}
