package postgres

import (
	"context"

	"github.com/steven230500/cartagena-api/internal/db"
	"github.com/steven230500/cartagena-api/internal/domain"
)

type EventRepository struct {
	Q *db.Queries
}

func NewEventRepository(q *db.Queries) *EventRepository {
	return &EventRepository{Q: q}
}

func toDomainEvent(e db.Event) domain.Event {
	return domain.Event{
		ID: e.ID.String(), Title: e.Title, Slug: e.Slug, StartDate: dateToString(e.StartDate),
		EndDate: dateToPtr(e.EndDate), Type: e.Type, Venue: e.Venue,
		RelatedBusinessID: uuidToPtr(e.RelatedBusinessID), Lat: e.Lat, Lng: e.Lng, Image: e.Image,
		Tags: e.Tags, Description: e.Description, Content: e.Content,
		CreatedAt: e.CreatedAt.Time, UpdatedAt: e.UpdatedAt.Time,
	}
}

func (r *EventRepository) List(ctx context.Context) ([]domain.Event, error) {
	rows, err := r.Q.ListEvents(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]domain.Event, len(rows))
	for i, e := range rows {
		out[i] = toDomainEvent(e)
	}
	return out, nil
}

func (r *EventRepository) GetBySlug(ctx context.Context, slug string) (domain.Event, error) {
	e, err := r.Q.GetEventBySlug(ctx, slug)
	if err != nil {
		return domain.Event{}, domain.ErrNotFound
	}
	return toDomainEvent(e), nil
}

func (r *EventRepository) Create(ctx context.Context, e domain.Event) (domain.Event, error) {
	startDate, err := parseDate(e.StartDate)
	if err != nil {
		return domain.Event{}, err
	}
	var endDateStr string
	if e.EndDate != nil {
		endDateStr = *e.EndDate
	}
	endDate, err := parseDate(endDateStr)
	if err != nil {
		return domain.Event{}, err
	}

	row, err := r.Q.CreateEvent(ctx, db.CreateEventParams{
		Title: e.Title, Slug: e.Slug, StartDate: startDate, EndDate: endDate,
		Type: e.Type, Venue: e.Venue, Lat: e.Lat, Lng: e.Lng, Image: e.Image,
		Tags: normalizeTags(e.Tags), Description: e.Description, Content: e.Content,
	})
	if err != nil {
		return domain.Event{}, err
	}
	return toDomainEvent(row), nil
}

func (r *EventRepository) Update(ctx context.Context, id string, e domain.Event) (domain.Event, error) {
	uid, err := mustUUID(id)
	if err != nil {
		return domain.Event{}, err
	}
	startDate, err := parseDate(e.StartDate)
	if err != nil {
		return domain.Event{}, err
	}
	var endDateStr string
	if e.EndDate != nil {
		endDateStr = *e.EndDate
	}
	endDate, err := parseDate(endDateStr)
	if err != nil {
		return domain.Event{}, err
	}

	row, err := r.Q.UpdateEvent(ctx, db.UpdateEventParams{
		ID: uid, Title: e.Title, Slug: e.Slug, StartDate: startDate, EndDate: endDate,
		Type: e.Type, Venue: e.Venue, Lat: e.Lat, Lng: e.Lng, Image: e.Image,
		Tags: normalizeTags(e.Tags), Description: e.Description, Content: e.Content,
	})
	if err != nil {
		return domain.Event{}, domain.ErrNotFound
	}
	return toDomainEvent(row), nil
}

func (r *EventRepository) Delete(ctx context.Context, id string) error {
	uid, err := mustUUID(id)
	if err != nil {
		return err
	}
	return r.Q.DeleteEvent(ctx, uid)
}
