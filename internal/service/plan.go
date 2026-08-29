package service

import (
	"context"

	"github.com/steven230500/cartagena-api/internal/domain"
	"github.com/steven230500/cartagena-api/internal/repository"
)

type PlanService struct {
	repo repository.PlanRepository
}

func NewPlanService(repo repository.PlanRepository) *PlanService {
	return &PlanService{repo: repo}
}

func (s *PlanService) List(ctx context.Context) ([]domain.Plan, error) {
	return s.repo.List(ctx)
}

func (s *PlanService) Create(ctx context.Context, p domain.Plan) (domain.Plan, error) {
	return s.repo.Create(ctx, p)
}

func (s *PlanService) Update(ctx context.Context, id string, p domain.Plan) (domain.Plan, error) {
	return s.repo.Update(ctx, id, p)
}

func (s *PlanService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
