package service

import (
	"context"

	"github.com/steven230500/cartagena-api/internal/domain"
	"github.com/steven230500/cartagena-api/internal/repository"
)

type BusinessService struct {
	repo repository.BusinessRepository
}

func NewBusinessService(repo repository.BusinessRepository) *BusinessService {
	return &BusinessService{repo: repo}
}

func (s *BusinessService) List(ctx context.Context, filter domain.BusinessFilter) ([]domain.Business, error) {
	return s.repo.List(ctx, filter)
}

func (s *BusinessService) GetBySlug(ctx context.Context, slug string) (domain.Business, error) {
	return s.repo.GetBySlug(ctx, slug)
}

func (s *BusinessService) Create(ctx context.Context, b domain.Business) (domain.Business, error) {
	return s.repo.Create(ctx, b)
}

func (s *BusinessService) Update(ctx context.Context, id string, b domain.Business) (domain.Business, error) {
	return s.repo.Update(ctx, id, b)
}

func (s *BusinessService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *BusinessService) ListByOwner(ctx context.Context, ownerID string) ([]domain.Business, error) {
	return s.repo.ListByOwner(ctx, ownerID)
}

// UpdateAsOwner solo deja tocar los campos de domain.BusinessOwnerPatch — nombre/
// slug/tipo/subtipo/barrio/coords quedan intactos pase lo que pase en el patch.
func (s *BusinessService) UpdateAsOwner(ctx context.Context, id, ownerID string, patch domain.BusinessOwnerPatch) (domain.Business, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Business{}, err
	}
	if existing.OwnerID == nil || *existing.OwnerID != ownerID {
		return domain.Business{}, domain.ErrForbidden
	}

	existing.Description = patch.Description
	existing.Hours = patch.Hours
	existing.PriceHint = patch.PriceHint
	existing.PriceTypicalNote = patch.PriceTypicalNote
	existing.Phone = patch.Phone
	existing.Web = patch.Web
	existing.Email = patch.Email
	existing.Instagram = patch.Instagram
	existing.Image = patch.Image
	existing.Tags = patch.Tags

	return s.repo.Update(ctx, id, existing)
}
