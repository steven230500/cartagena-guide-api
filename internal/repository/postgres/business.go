package postgres

import (
	"context"

	"github.com/steven230500/cartagena-api/internal/db"
	"github.com/steven230500/cartagena-api/internal/domain"
)

type BusinessRepository struct {
	Q *db.Queries
}

func NewBusinessRepository(q *db.Queries) *BusinessRepository {
	return &BusinessRepository{Q: q}
}

func toDomainBusiness(b db.Business) domain.Business {
	return domain.Business{
		ID: b.ID.String(), Name: b.Name, Slug: b.Slug, Type: b.Type, Subtype: b.Subtype,
		Barrio: b.Barrio, Lat: b.Lat, Lng: b.Lng, Image: b.Image, Tags: b.Tags,
		Description: b.Description, DescriptionEn: b.DescriptionEn, Hours: b.Hours, PriceHint: b.PriceHint,
		PriceTypicalNote: b.PriceTypicalNote, Phone: b.Phone, Web: b.Web, Email: b.Email,
		Instagram: b.Instagram, OwnerID: uuidToPtr(b.OwnerID), Verified: b.Verified,
		CreatedAt: b.CreatedAt.Time, UpdatedAt: b.UpdatedAt.Time,
	}
}

func (r *BusinessRepository) List(ctx context.Context, filter domain.BusinessFilter) ([]domain.Business, error) {
	rows, err := r.Q.ListBusinesses(ctx, db.ListBusinessesParams{
		Type:   strOrNil(filter.Type),
		Barrio: strOrNil(filter.Barrio),
		Q:      strOrNil(filter.Q),
	})
	if err != nil {
		return nil, err
	}
	out := make([]domain.Business, len(rows))
	for i, b := range rows {
		out[i] = toDomainBusiness(b)
	}
	return out, nil
}

func (r *BusinessRepository) GetBySlug(ctx context.Context, slug string) (domain.Business, error) {
	b, err := r.Q.GetBusinessBySlug(ctx, slug)
	if err != nil {
		return domain.Business{}, domain.ErrNotFound
	}
	return toDomainBusiness(b), nil
}

func (r *BusinessRepository) GetByID(ctx context.Context, id string) (domain.Business, error) {
	uid, err := mustUUID(id)
	if err != nil {
		return domain.Business{}, err
	}
	b, err := r.Q.GetBusinessByID(ctx, uid)
	if err != nil {
		return domain.Business{}, domain.ErrNotFound
	}
	return toDomainBusiness(b), nil
}

func (r *BusinessRepository) ListByOwner(ctx context.Context, ownerID string) ([]domain.Business, error) {
	uid, err := mustUUID(ownerID)
	if err != nil {
		return nil, err
	}
	rows, err := r.Q.ListBusinessesByOwner(ctx, uid)
	if err != nil {
		return nil, err
	}
	out := make([]domain.Business, len(rows))
	for i, b := range rows {
		out[i] = toDomainBusiness(b)
	}
	return out, nil
}

func (r *BusinessRepository) SetOwner(ctx context.Context, id, ownerID string) (domain.Business, error) {
	uid, err := mustUUID(id)
	if err != nil {
		return domain.Business{}, err
	}
	oid, err := mustUUID(ownerID)
	if err != nil {
		return domain.Business{}, err
	}
	b, err := r.Q.SetBusinessOwner(ctx, db.SetBusinessOwnerParams{ID: uid, OwnerID: oid})
	if err != nil {
		return domain.Business{}, domain.ErrNotFound
	}
	return toDomainBusiness(b), nil
}

func (r *BusinessRepository) Create(ctx context.Context, b domain.Business) (domain.Business, error) {
	row, err := r.Q.CreateBusiness(ctx, db.CreateBusinessParams{
		Name: b.Name, Slug: b.Slug, Type: b.Type, Subtype: b.Subtype, Barrio: b.Barrio,
		Lat: b.Lat, Lng: b.Lng, Image: b.Image, Tags: normalizeTags(b.Tags),
		Description: b.Description, DescriptionEn: b.DescriptionEn, Hours: b.Hours, PriceHint: b.PriceHint,
		PriceTypicalNote: b.PriceTypicalNote, Phone: b.Phone, Web: b.Web, Email: b.Email,
		Instagram: b.Instagram, Verified: b.Verified,
	})
	if err != nil {
		return domain.Business{}, err
	}
	return toDomainBusiness(row), nil
}

func (r *BusinessRepository) Update(ctx context.Context, id string, b domain.Business) (domain.Business, error) {
	uid, err := mustUUID(id)
	if err != nil {
		return domain.Business{}, err
	}
	row, err := r.Q.UpdateBusiness(ctx, db.UpdateBusinessParams{
		ID: uid, Name: b.Name, Slug: b.Slug, Type: b.Type, Subtype: b.Subtype, Barrio: b.Barrio,
		Lat: b.Lat, Lng: b.Lng, Image: b.Image, Tags: normalizeTags(b.Tags),
		Description: b.Description, DescriptionEn: b.DescriptionEn, Hours: b.Hours, PriceHint: b.PriceHint,
		PriceTypicalNote: b.PriceTypicalNote, Phone: b.Phone, Web: b.Web, Email: b.Email,
		Instagram: b.Instagram, Verified: b.Verified,
	})
	if err != nil {
		return domain.Business{}, domain.ErrNotFound
	}
	return toDomainBusiness(row), nil
}

func (r *BusinessRepository) Delete(ctx context.Context, id string) error {
	uid, err := mustUUID(id)
	if err != nil {
		return err
	}
	return r.Q.DeleteBusiness(ctx, uid)
}

// strOrNil convierte un filtro vacío en nil para los sqlc.narg opcionales de ListBusinesses.
func strOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
