package domain

import "time"

// StartDate/EndDate quedan como string "2006-01-02" — mismo formato que la API
// ya devuelve hoy, para no romper el contrato con el front.
type Event struct {
	ID                string    `json:"id"`
	Title             string    `json:"title"`
	Slug              string    `json:"slug"`
	StartDate         string    `json:"start_date"`
	EndDate           *string   `json:"end_date"`
	Type              string    `json:"type"`
	Venue             string    `json:"venue"`
	RelatedBusinessID *string   `json:"related_business_id"`
	Lat               float64   `json:"lat"`
	Lng               float64   `json:"lng"`
	Image             string    `json:"image"`
	Tags              []string  `json:"tags"`
	Description       string    `json:"description"`
	Content           *string   `json:"content"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
