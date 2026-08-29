package service

import (
	"context"

	"github.com/steven230500/cartagena-api/internal/domain"
	"github.com/steven230500/cartagena-api/internal/repository"
)

type FavoriteService struct {
	repo repository.FavoriteRepository
}

func NewFavoriteService(repo repository.FavoriteRepository) *FavoriteService {
	return &FavoriteService{repo: repo}
}

func (s *FavoriteService) List(ctx context.Context, userID string) ([]domain.Business, error) {
	return s.repo.List(ctx, userID)
}

func (s *FavoriteService) Add(ctx context.Context, userID, businessID string) error {
	return s.repo.Add(ctx, userID, businessID)
}

func (s *FavoriteService) Remove(ctx context.Context, userID, businessID string) error {
	return s.repo.Remove(ctx, userID, businessID)
}
