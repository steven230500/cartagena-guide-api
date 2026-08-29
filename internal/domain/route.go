package domain

import "time"

type RouteStep struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	AudioURL    *string  `json:"audio_url"`
	Duration    *string  `json:"duration"`
	Lat         *float64 `json:"lat"`
	Lng         *float64 `json:"lng"`
	Image       *string  `json:"image"`
	Position    int32    `json:"position"`
}

type Route struct {
	ID          string      `json:"id"`
	Slug        string      `json:"slug"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Duration    string      `json:"duration"`
	Distance    string      `json:"distance"`
	Difficulty  string      `json:"difficulty"`
	Category    string      `json:"category"`
	Image       string      `json:"image"`
	Highlights  []string    `json:"highlights"`
	AudioGuide  bool        `json:"audio_guide"`
	Offline     bool        `json:"offline"`
	Price       string      `json:"price"`
	Steps       []RouteStep `json:"steps"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}
