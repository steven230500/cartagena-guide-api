package service

import (
	"context"

	"github.com/steven230500/cartagena-api/internal/domain"
	"github.com/steven230500/cartagena-api/internal/repository"
)

type ParishService struct {
	repo repository.ParishRepository
}

func NewParishService(repo repository.ParishRepository) *ParishService {
	return &ParishService{repo: repo}
}

func (s *ParishService) List(ctx context.Context) ([]domain.Parish, error) {
	return s.repo.List(ctx)
}

func (s *ParishService) GetByID(ctx context.Context, id string) (domain.Parish, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ParishService) Create(ctx context.Context, p domain.Parish) (domain.Parish, error) {
	return s.repo.Create(ctx, p)
}

func (s *ParishService) Update(ctx context.Context, id string, p domain.Parish) (domain.Parish, error) {
	return s.repo.Update(ctx, id, p)
}

func (s *ParishService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
