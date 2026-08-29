package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/steven230500/cartagena-api/internal/db"
	"github.com/steven230500/cartagena-api/internal/domain"
)

type ParishRepository struct {
	Q *db.Queries
}

func NewParishRepository(q *db.Queries) *ParishRepository {
	return &ParishRepository{Q: q}
}

func toDomainSchedule(s db.ParishSchedule) domain.Schedule {
	return domain.Schedule{
		ID: s.ID.String(), ParishID: s.ParishID.String(), Day: s.Day, Times: s.Times, Position: s.Position,
	}
}

func (r *ParishRepository) withSchedules(ctx context.Context, p db.Parish) (domain.Parish, error) {
	rows, err := r.Q.ListSchedulesByParish(ctx, p.ID)
	if err != nil {
		return domain.Parish{}, err
	}
	schedules := make([]domain.Schedule, len(rows))
	for i, s := range rows {
		schedules[i] = toDomainSchedule(s)
	}
	return domain.Parish{
		ID: p.ID.String(), Name: p.Name, Address: p.Address, Neighborhood: p.Neighborhood,
		Phone: p.Phone, Schedules: schedules,
	}, nil
}

// replaceSchedules borra todos los horarios de la parroquia y reinserta los nuevos.
// Mismo comportamiento que ya tenía el handler antes del refactor.
func (r *ParishRepository) replaceSchedules(ctx context.Context, parishID pgtype.UUID, schedules []domain.Schedule) error {
	if err := r.Q.DeleteSchedulesByParish(ctx, parishID); err != nil {
		return err
	}
	for i, s := range schedules {
		times := s.Times
		if times == nil {
			times = []string{}
		}
		if _, err := r.Q.CreateParishSchedule(ctx, db.CreateParishScheduleParams{
			ParishID: parishID, Day: s.Day, Times: times, Position: int32(i),
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *ParishRepository) List(ctx context.Context) ([]domain.Parish, error) {
	rows, err := r.Q.ListParishes(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]domain.Parish, len(rows))
	for i, p := range rows {
		full, err := r.withSchedules(ctx, p)
		if err != nil {
			return nil, err
		}
		out[i] = full
	}
	return out, nil
}

func (r *ParishRepository) GetByID(ctx context.Context, id string) (domain.Parish, error) {
	parishes, err := r.List(ctx)
	if err != nil {
		return domain.Parish{}, err
	}
	for _, p := range parishes {
		if p.ID == id {
			return p, nil
		}
	}
	return domain.Parish{}, domain.ErrNotFound
}

func (r *ParishRepository) Create(ctx context.Context, p domain.Parish) (domain.Parish, error) {
	row, err := r.Q.CreateParish(ctx, db.CreateParishParams{
		Name: p.Name, Address: p.Address, Neighborhood: p.Neighborhood, Phone: p.Phone,
	})
	if err != nil {
		return domain.Parish{}, err
	}
	if err := r.replaceSchedules(ctx, row.ID, p.Schedules); err != nil {
		return domain.Parish{}, err
	}
	return r.withSchedules(ctx, row)
}

func (r *ParishRepository) Update(ctx context.Context, id string, p domain.Parish) (domain.Parish, error) {
	uid, err := mustUUID(id)
	if err != nil {
		return domain.Parish{}, err
	}
	row, err := r.Q.UpdateParish(ctx, db.UpdateParishParams{
		ID: uid, Name: p.Name, Address: p.Address, Neighborhood: p.Neighborhood, Phone: p.Phone,
	})
	if err != nil {
		return domain.Parish{}, domain.ErrNotFound
	}
	if err := r.replaceSchedules(ctx, uid, p.Schedules); err != nil {
		return domain.Parish{}, err
	}
	return r.withSchedules(ctx, row)
}

func (r *ParishRepository) Delete(ctx context.Context, id string) error {
	uid, err := mustUUID(id)
	if err != nil {
		return err
	}
	// parish_schedules tiene ON DELETE CASCADE hacia parishes.
	return r.Q.DeleteParish(ctx, uid)
}
