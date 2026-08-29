package domain

type Schedule struct {
	ID       string   `json:"id"`
	ParishID string   `json:"parish_id"`
	Day      string   `json:"day"`
	Times    []string `json:"times"`
	Position int32    `json:"position"`
}

type Parish struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Address      string     `json:"address"`
	Neighborhood string     `json:"neighborhood"`
	Phone        *string    `json:"phone"`
	Schedules    []Schedule `json:"schedules"`
}
