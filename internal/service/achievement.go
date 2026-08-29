package service

import (
	"context"

	"github.com/steven230500/cartagena-api/internal/domain"
	"github.com/steven230500/cartagena-api/internal/repository"
)

type AchievementService struct {
	repo              repository.AchievementRepository
	favoriteRepo      repository.FavoriteRepository
	routeProgressRepo repository.RouteProgressRepository
}

func NewAchievementService(
	repo repository.AchievementRepository,
	favoriteRepo repository.FavoriteRepository,
	routeProgressRepo repository.RouteProgressRepository,
) *AchievementService {
	return &AchievementService{repo: repo, favoriteRepo: favoriteRepo, routeProgressRepo: routeProgressRepo}
}

func (s *AchievementService) List(ctx context.Context) ([]domain.Achievement, error) {
	return s.repo.List(ctx)
}

func (s *AchievementService) Create(ctx context.Context, a domain.Achievement) (domain.Achievement, error) {
	return s.repo.Create(ctx, a)
}

func (s *AchievementService) Update(ctx context.Context, id string, a domain.Achievement) (domain.Achievement, error) {
	return s.repo.Update(ctx, id, a)
}

func (s *AchievementService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *AchievementService) Progress(ctx context.Context, userID string) ([]domain.AchievementProgress, domain.AchievementStats, error) {
	defs, err := s.repo.List(ctx)
	if err != nil {
		return nil, domain.AchievementStats{}, err
	}

	favorites, err := s.favoriteRepo.List(ctx, userID)
	if err != nil {
		return nil, domain.AchievementStats{}, err
	}
	routesCompleted, err := s.routeProgressRepo.CountCompleted(ctx, userID)
	if err != nil {
		return nil, domain.AchievementStats{}, err
	}

	stats := domain.AchievementStats{
		RoutesCompleted: routesCompleted,
		FavoritesCount:  len(favorites),
	}

	out := make([]domain.AchievementProgress, len(defs))
	for i, def := range defs {
		current := 0
		switch def.CriteriaType {
		case "routes_completed":
			current = stats.RoutesCompleted
		case "favorites_count":
			current = stats.FavoritesCount
		}
		out[i] = domain.AchievementProgress{
			Achievement: def,
			Current:     current,
			Unlocked:    current >= def.Threshold,
		}
	}

	return out, stats, nil
}
