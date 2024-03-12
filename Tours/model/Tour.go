package model

import (
	"time"
)

type DifficultyLevel int

const (
	Easy DifficultyLevel = iota
	Moderate
	Difficult
)

type TourStatus int

const (
	Draft TourStatus = iota
	Published
	Archived
)

type Tour struct {
	Name              string          `json:"name"`
	DifficultyLevel   DifficultyLevel `json:"difficultyLevel"`
	Description       string          `json:"description"`
	Status            TourStatus      `json:"status"`
	Price             int             `json:"price"`
	UserId            int             `json:"userId"`
	PublishedDateTime *time.Time      `json:"publishedDateTime,omitempty"`
	ArchivedDateTime  *time.Time      `json:"archivedDateTime,omitempty"`
}
