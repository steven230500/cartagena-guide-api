package service

import (
	"context"

	"github.com/steven230500/cartagena-api/internal/domain"
	"github.com/steven230500/cartagena-api/internal/repository"
)

type RouteProgressService struct {
	repo repository.RouteProgressRepository
}

func NewRouteProgressService(repo repository.RouteProgressRepository) *RouteProgressService {
	return &RouteProgressService{repo: repo}
}

func (s *RouteProgressService) Get(ctx context.Context, userID, routeID string) (domain.RouteProgress, error) {
	return s.repo.Get(ctx, userID, routeID)
}

func (s *RouteProgressService) Upsert(ctx context.Context, userID, routeID string, currentStep int) (domain.RouteProgress, error) {
	return s.repo.Upsert(ctx, userID, routeID, currentStep)
}
