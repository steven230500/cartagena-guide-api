package domain

import "time"

type Business struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Slug             string    `json:"slug"`
	Type             string    `json:"type"`
	Subtype          string    `json:"subtype"`
	Barrio           string    `json:"barrio"`
	Lat              float64   `json:"lat"`
	Lng              float64   `json:"lng"`
	Image            string    `json:"image"`
	Tags             []string  `json:"tags"`
	Description      string    `json:"description"`
	DescriptionEn    *string   `json:"description_en"`
	Hours            *string   `json:"hours"`
	PriceHint        *string   `json:"price_hint"`
	PriceTypicalNote *string   `json:"price_typical_note"`
	Phone            *string   `json:"phone"`
	Web              *string   `json:"web"`
	Email            *string   `json:"email"`
	Instagram        *string   `json:"instagram"`
	OwnerID          *string   `json:"owner_id"`
	Verified         bool      `json:"verified"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// BusinessFilter son los filtros opcionales de GET /api/businesses.
type BusinessFilter struct {
	Type   string
	Barrio string
	Q      string
}

// BusinessOwnerPatch es el subconjunto de campos que puede tocar el dueño
// de un negocio (Fase 5) — nombre/slug/tipo/subtipo/barrio/coords/verified
// quedan admin-only, se pisan siempre con los valores existentes.
type BusinessOwnerPatch struct {
	Description      string
	Hours            *string
	PriceHint        *string
	PriceTypicalNote *string
	Phone            *string
	Web              *string
	Email            *string
	Instagram        *string
	Image            string
	Tags             []string
}
