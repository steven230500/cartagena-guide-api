package service

import (
	"context"

	"github.com/steven230500/cartagena-api/internal/domain"
	"github.com/steven230500/cartagena-api/internal/repository"
)

type EventService struct {
	repo repository.EventRepository
}

func NewEventService(repo repository.EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) List(ctx context.Context) ([]domain.Event, error) {
	return s.repo.List(ctx)
}

func (s *EventService) GetBySlug(ctx context.Context, slug string) (domain.Event, error) {
	return s.repo.GetBySlug(ctx, slug)
}

func (s *EventService) Create(ctx context.Context, e domain.Event) (domain.Event, error) {
	return s.repo.Create(ctx, e)
}

func (s *EventService) Update(ctx context.Context, id string, e domain.Event) (domain.Event, error) {
	return s.repo.Update(ctx, id, e)
}

func (s *EventService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
