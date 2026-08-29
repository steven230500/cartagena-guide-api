package postgres

import (
	"context"

	"github.com/steven230500/cartagena-api/internal/db"
	"github.com/steven230500/cartagena-api/internal/domain"
)

type PlanRepository struct {
	Q *db.Queries
}

func NewPlanRepository(q *db.Queries) *PlanRepository {
	return &PlanRepository{Q: q}
}

func toDomainPlan(p db.Plan) domain.Plan {
	return domain.Plan{
		ID: p.ID.String(), Title: p.Title, Description: p.Description, Type: p.Type,
		Price: p.Price, PriceAmount: p.PriceAmount, Date: p.Date, Time: p.Time,
		Location: p.Location, Neighborhood: p.Neighborhood,
	}
}

func (r *PlanRepository) List(ctx context.Context) ([]domain.Plan, error) {
	rows, err := r.Q.ListPlans(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]domain.Plan, len(rows))
	for i, p := range rows {
		out[i] = toDomainPlan(p)
	}
	return out, nil
}

func (r *PlanRepository) Create(ctx context.Context, p domain.Plan) (domain.Plan, error) {
	row, err := r.Q.CreatePlan(ctx, db.CreatePlanParams{
		Title: p.Title, Description: p.Description, Type: p.Type, Price: p.Price,
		PriceAmount: p.PriceAmount, Date: p.Date, Time: p.Time, Location: p.Location,
		Neighborhood: p.Neighborhood,
	})
	if err != nil {
		return domain.Plan{}, err
	}
	return toDomainPlan(row), nil
}

func (r *PlanRepository) Update(ctx context.Context, id string, p domain.Plan) (domain.Plan, error) {
	uid, err := mustUUID(id)
	if err != nil {
		return domain.Plan{}, err
	}
	row, err := r.Q.UpdatePlan(ctx, db.UpdatePlanParams{
		ID: uid, Title: p.Title, Description: p.Description, Type: p.Type, Price: p.Price,
		PriceAmount: p.PriceAmount, Date: p.Date, Time: p.Time, Location: p.Location,
		Neighborhood: p.Neighborhood,
	})
	if err != nil {
		return domain.Plan{}, domain.ErrNotFound
	}
	return toDomainPlan(row), nil
}

func (r *PlanRepository) Delete(ctx context.Context, id string) error {
	uid, err := mustUUID(id)
	if err != nil {
		return err
	}
	return r.Q.DeletePlan(ctx, uid)
}
