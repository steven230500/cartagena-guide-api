package postgres

import (
	"context"

	"github.com/steven230500/cartagena-api/internal/db"
	"github.com/steven230500/cartagena-api/internal/domain"
)

type AchievementRepository struct {
	Q *db.Queries
}

func NewAchievementRepository(q *db.Queries) *AchievementRepository {
	return &AchievementRepository{Q: q}
}

func toDomainAchievement(a db.Achievement) domain.Achievement {
	return domain.Achievement{
		ID:           a.ID.String(),
		Code:         a.Code,
		Title:        a.Title,
		Description:  a.Description,
		Icon:         a.Icon,
		CriteriaType: a.CriteriaType,
		Threshold:    int(a.Threshold),
	}
}

func (r *AchievementRepository) List(ctx context.Context) ([]domain.Achievement, error) {
	rows, err := r.Q.ListAchievements(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]domain.Achievement, len(rows))
	for i, a := range rows {
		out[i] = toDomainAchievement(a)
	}
	return out, nil
}

func (r *AchievementRepository) Create(ctx context.Context, in domain.Achievement) (domain.Achievement, error) {
	a, err := r.Q.CreateAchievement(ctx, db.CreateAchievementParams{
		Code: in.Code, Title: in.Title, Description: in.Description,
		Icon: in.Icon, CriteriaType: in.CriteriaType, Threshold: int32(in.Threshold),
	})
	if err != nil {
		return domain.Achievement{}, err
	}
	return toDomainAchievement(a), nil
}

func (r *AchievementRepository) Update(ctx context.Context, id string, in domain.Achievement) (domain.Achievement, error) {
	uid, err := mustUUID(id)
	if err != nil {
		return domain.Achievement{}, err
	}
	a, err := r.Q.UpdateAchievement(ctx, db.UpdateAchievementParams{
		ID: uid, Code: in.Code, Title: in.Title, Description: in.Description,
		Icon: in.Icon, CriteriaType: in.CriteriaType, Threshold: int32(in.Threshold),
	})
	if err != nil {
		return domain.Achievement{}, domain.ErrNotFound
	}
	return toDomainAchievement(a), nil
}

func (r *AchievementRepository) Delete(ctx context.Context, id string) error {
	uid, err := mustUUID(id)
	if err != nil {
		return err
	}
	return r.Q.DeleteAchievement(ctx, uid)
}
