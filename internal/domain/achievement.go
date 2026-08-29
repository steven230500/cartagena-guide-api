package domain

type Achievement struct {
	ID           string `json:"id"`
	Code         string `json:"code"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Icon         string `json:"icon"`
	CriteriaType string `json:"criteria_type"`
	Threshold    int    `json:"threshold"`
}

type AchievementProgress struct {
	Achievement
	Current  int  `json:"current"`
	Unlocked bool `json:"unlocked"`
}

type AchievementStats struct {
	RoutesCompleted int `json:"routes_completed"`
	FavoritesCount  int `json:"favorites_count"`
}
