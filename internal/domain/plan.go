package domain

type Plan struct {
	ID           string  `json:"id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Type         string  `json:"type"`
	Price        string  `json:"price"`
	PriceAmount  *string `json:"price_amount"`
	Date         string  `json:"date"`
	Time         string  `json:"time"`
	Location     string  `json:"location"`
	Neighborhood string  `json:"neighborhood"`
}
