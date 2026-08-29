package domain

type RouteProgress struct {
	RouteID     string `json:"route_id"`
	CurrentStep int    `json:"current_step"`
	UpdatedAt   string `json:"updated_at"`
}
