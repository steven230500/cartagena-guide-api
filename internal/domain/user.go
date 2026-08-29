package domain

import "time"

// User nunca lleva el password hash — ese vive solo dentro del repository,
// así un json.Marshal accidental de User jamás puede filtrarlo.
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
