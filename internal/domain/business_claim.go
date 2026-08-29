package domain

import "time"

type BusinessClaim struct {
	ID           string     `json:"id"`
	BusinessID   string     `json:"business_id"`
	BusinessName string     `json:"business_name"`
	BusinessSlug string     `json:"business_slug"`
	UserID       string     `json:"user_id"`
	UserEmail    string     `json:"user_email"`
	Message      string     `json:"message"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	ResolvedAt   *time.Time `json:"resolved_at"`
}

type BusinessClaimFilter struct {
	Status string
	UserID string
}
