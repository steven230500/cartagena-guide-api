package service

import (
	"context"
	"errors"

	"github.com/steven230500/cartagena-api/internal/domain"
	"github.com/steven230500/cartagena-api/internal/repository"
)

var ErrAlreadyOwned = errors.New("business already has an owner")

type BusinessClaimService struct {
	repo         repository.BusinessClaimRepository
	businessRepo repository.BusinessRepository
	userRepo     repository.UserRepository
}

func NewBusinessClaimService(
	repo repository.BusinessClaimRepository,
	businessRepo repository.BusinessRepository,
	userRepo repository.UserRepository,
) *BusinessClaimService {
	return &BusinessClaimService{repo: repo, businessRepo: businessRepo, userRepo: userRepo}
}

func (s *BusinessClaimService) Create(ctx context.Context, businessID, userID, message string) (domain.BusinessClaim, error) {
	business, err := s.businessRepo.GetByID(ctx, businessID)
	if err != nil {
		return domain.BusinessClaim{}, err
	}
	if business.OwnerID != nil {
		return domain.BusinessClaim{}, ErrAlreadyOwned
	}
	return s.repo.Create(ctx, businessID, userID, message)
}

func (s *BusinessClaimService) ListPending(ctx context.Context) ([]domain.BusinessClaim, error) {
	return s.repo.List(ctx, domain.BusinessClaimFilter{Status: "pending"})
}

func (s *BusinessClaimService) ListMine(ctx context.Context, userID string) ([]domain.BusinessClaim, error) {
	return s.repo.List(ctx, domain.BusinessClaimFilter{UserID: userID})
}

// Approve toca 3 tablas con 3 updates secuenciales, no una transacción real —
// mismo criterio no-transaccional que el resto del proyecto (borra-y-reinserta
// de route_steps/parish_schedules). Riesgo bajo: acción manual de admin, no concurrente.
func (s *BusinessClaimService) Approve(ctx context.Context, claimID string) (domain.BusinessClaim, error) {
	claim, err := s.repo.Get(ctx, claimID)
	if err != nil {
		return domain.BusinessClaim{}, err
	}

	resolved, err := s.repo.Resolve(ctx, claimID, "approved")
	if err != nil {
		return domain.BusinessClaim{}, err
	}
	if _, err := s.businessRepo.SetOwner(ctx, claim.BusinessID, claim.UserID); err != nil {
		return domain.BusinessClaim{}, err
	}
	if _, err := s.userRepo.UpdateRole(ctx, claim.UserID, "business_owner"); err != nil {
		return domain.BusinessClaim{}, err
	}

	return resolved, nil
}

func (s *BusinessClaimService) Reject(ctx context.Context, claimID string) (domain.BusinessClaim, error) {
	return s.repo.Resolve(ctx, claimID, "rejected")
}
