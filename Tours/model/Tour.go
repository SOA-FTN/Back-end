package model

import (
	"time"

	"gorm.io/gorm"
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
	gorm.Model
	Name              string          `json:"name"`
	DifficultyLevel   DifficultyLevel `json:"difficultyLevel"`
	Description       string          `json:"description"`
	TStatus           TourStatus      `json:"tStatus"`
	Price             int             `json:"price"`
	UserId            int             `json:"userId"`
	PublishedDateTime *time.Time      `json:"publishedDateTime,omitempty"`
	ArchivedDateTime  *time.Time      `json:"archivedDateTime,omitempty"`
}

type CreateTourRequest struct {
	ID                  uint          `json:"id"`
	Name                string        `json:"name"`
	Description         string        `json:"description"`
	DifficultyLevel     string        `json:"DifficultyLevel"`
	Status              string        `json:"status"`
	Price               int           `json:"price"`
	UserID              int           `json:"userId"`
	PublishedDateTime   *time.Time    `json:"publishedDateTime,omitempty"`
	ArchivedDateTime    *time.Time    `json:"archivedDateTime,omitempty"`
	Tags                []string      `json:"tags"`
	TourPoints          []interface{} `json:"tourPoints"`
	TourCharacteristics []interface{} `json:"tourCharacteristics"`
	TourReviews         []interface{} `json:"tourReviews"`
}
