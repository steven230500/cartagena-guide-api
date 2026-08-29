package service

import (
	"context"

	"github.com/steven230500/cartagena-api/internal/domain"
	"github.com/steven230500/cartagena-api/internal/repository"
)

type RouteService struct {
	repo repository.RouteRepository
}

func NewRouteService(repo repository.RouteRepository) *RouteService {
	return &RouteService{repo: repo}
}

func (s *RouteService) List(ctx context.Context) ([]domain.Route, error) {
	return s.repo.List(ctx)
}

func (s *RouteService) GetBySlug(ctx context.Context, slug string) (domain.Route, error) {
	return s.repo.GetBySlug(ctx, slug)
}

func (s *RouteService) Create(ctx context.Context, r domain.Route) (domain.Route, error) {
	return s.repo.Create(ctx, r)
}

func (s *RouteService) Update(ctx context.Context, id string, r domain.Route) (domain.Route, error) {
	return s.repo.Update(ctx, id, r)
}

func (s *RouteService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
