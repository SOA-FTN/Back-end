package model

type TourPoint struct {
	TourId      int64   `json:"tourId"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	ImageUrl    string  `json:"imageUrl"`
	Secret      string  `json:"secret"`
}
