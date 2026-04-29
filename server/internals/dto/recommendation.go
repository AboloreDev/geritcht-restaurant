package dto

import "time"

type RecommendationRequest struct {
	Allergies []string `json:"allergies"`
	Dietary   string   `json:"dietary"`
	Mood      string   `json:"mood" binding:"required,oneof=light heavy spicy surprise"`
	Budget    string   `json:"budget" binding:"required,oneof=under_3k 3k_to_8k no_limit"`
}

type RecommendedItemResponse struct {
	Item   MenuResponse `json:"item"`
	Reason string       `json:"reason"`
	Score  float64      `json:"score"`
}

type RecommendationResponse struct {
	ID              uint                      `json:"id"`
	Recommendations []RecommendedItemResponse `json:"recommendations"`
	Preferences     RecommendationRequest     `json:"preferences"`
	CreatedAt       time.Time                 `json:"created_at"`
}
