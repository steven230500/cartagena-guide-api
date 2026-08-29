package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/steven230500/cartagena-api/internal/db"
	"github.com/steven230500/cartagena-api/internal/domain"
)

type RouteRepository struct {
	Q *db.Queries
}

func NewRouteRepository(q *db.Queries) *RouteRepository {
	return &RouteRepository{Q: q}
}

func toDomainStep(s db.RouteStep) domain.RouteStep {
	return domain.RouteStep{
		ID: s.ID.String(), Title: s.Title, Description: s.Description,
		AudioURL: s.AudioUrl, Duration: s.Duration, Lat: s.Lat, Lng: s.Lng,
		Image: s.Image, Position: s.Position,
	}
}

func (r *RouteRepository) withSteps(ctx context.Context, route db.Route) (domain.Route, error) {
	rows, err := r.Q.ListStepsByRoute(ctx, route.ID)
	if err != nil {
		return domain.Route{}, err
	}
	steps := make([]domain.RouteStep, len(rows))
	for i, s := range rows {
		steps[i] = toDomainStep(s)
	}
	return domain.Route{
		ID: route.ID.String(), Slug: route.Slug, Title: route.Title, Description: route.Description,
		Duration: route.Duration, Distance: route.Distance, Difficulty: route.Difficulty,
		Category: route.Category, Image: route.Image, Highlights: route.Highlights,
		AudioGuide: route.AudioGuide, Offline: route.Offline, Price: route.Price, Steps: steps,
		CreatedAt: route.CreatedAt.Time, UpdatedAt: route.UpdatedAt.Time,
	}, nil
}

// replaceSteps borra todos los pasos de la ruta y reinserta los nuevos —
// mismo patrón que ParishRepository.replaceSchedules.
func (r *RouteRepository) replaceSteps(ctx context.Context, routeID pgtype.UUID, steps []domain.RouteStep) error {
	if err := r.Q.DeleteStepsByRoute(ctx, routeID); err != nil {
		return err
	}
	for i, s := range steps {
		if _, err := r.Q.CreateRouteStep(ctx, db.CreateRouteStepParams{
			RouteID: routeID, Title: s.Title, Description: s.Description,
			AudioUrl: s.AudioURL, Duration: s.Duration, Lat: s.Lat, Lng: s.Lng,
			Image: s.Image, Position: int32(i),
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *RouteRepository) List(ctx context.Context) ([]domain.Route, error) {
	rows, err := r.Q.ListRoutes(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]domain.Route, len(rows))
	for i, route := range rows {
		full, err := r.withSteps(ctx, route)
		if err != nil {
			return nil, err
		}
		out[i] = full
	}
	return out, nil
}

func (r *RouteRepository) GetBySlug(ctx context.Context, slug string) (domain.Route, error) {
	route, err := r.Q.GetRouteBySlug(ctx, slug)
	if err != nil {
		return domain.Route{}, domain.ErrNotFound
	}
	return r.withSteps(ctx, route)
}

func (r *RouteRepository) Create(ctx context.Context, in domain.Route) (domain.Route, error) {
	route, err := r.Q.CreateRoute(ctx, db.CreateRouteParams{
		Slug: in.Slug, Title: in.Title, Description: in.Description, Duration: in.Duration,
		Distance: in.Distance, Difficulty: in.Difficulty, Category: in.Category, Image: in.Image,
		Highlights: normalizeTags(in.Highlights), AudioGuide: in.AudioGuide, Offline: in.Offline,
		Price: in.Price,
	})
	if err != nil {
		return domain.Route{}, err
	}
	if err := r.replaceSteps(ctx, route.ID, in.Steps); err != nil {
		return domain.Route{}, err
	}
	return r.withSteps(ctx, route)
}

func (r *RouteRepository) Update(ctx context.Context, id string, in domain.Route) (domain.Route, error) {
	uid, err := mustUUID(id)
	if err != nil {
		return domain.Route{}, err
	}
	route, err := r.Q.UpdateRoute(ctx, db.UpdateRouteParams{
		ID: uid, Slug: in.Slug, Title: in.Title, Description: in.Description, Duration: in.Duration,
		Distance: in.Distance, Difficulty: in.Difficulty, Category: in.Category, Image: in.Image,
		Highlights: normalizeTags(in.Highlights), AudioGuide: in.AudioGuide, Offline: in.Offline,
		Price: in.Price,
	})
	if err != nil {
		return domain.Route{}, domain.ErrNotFound
	}
	if err := r.replaceSteps(ctx, uid, in.Steps); err != nil {
		return domain.Route{}, err
	}
	return r.withSteps(ctx, route)
}

func (r *RouteRepository) Delete(ctx context.Context, id string) error {
	uid, err := mustUUID(id)
	if err != nil {
		return err
	}
	// route_steps tiene ON DELETE CASCADE hacia routes.
	return r.Q.DeleteRoute(ctx, uid)
}
