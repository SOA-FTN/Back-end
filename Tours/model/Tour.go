package model

import (
	"time"
	TOURS "tours/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
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
	Status           TourStatus      `json:"status"`
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

type PurchasedTours struct{
	UserId            int             `json:"userId"`
	TourId            int             `json:"tourId"`

}

func ToModelDifficultyLevel(difficultyLevel TOURS.DifficultyLevel) DifficultyLevel {
    switch difficultyLevel {
    case TOURS.DifficultyLevel_Easy:
        return Easy
    case TOURS.DifficultyLevel_Moderate:
        return Moderate
    case TOURS.DifficultyLevel_Difficult:
        return Difficult
    default:
        return Easy // or handle default case appropriately
    }
}

func ToModelStatus(status TOURS.TourStatus) TourStatus {
	switch status{
	case TOURS.TourStatus_Draft:
		return Draft
	case TOURS.TourStatus_Published:
		return Published
	case TOURS.TourStatus_Archived:
		return Archived
	default:
		return Draft
	}
}

func ToProtoDifficultyLevel(difficultyLevel DifficultyLevel) TOURS.DifficultyLevel {
	switch difficultyLevel {
	case Easy:
		return TOURS.DifficultyLevel_Easy
	case Moderate:
		return TOURS.DifficultyLevel_Moderate
	case Difficult:
		return TOURS.DifficultyLevel_Difficult
	default:
		return TOURS.DifficultyLevel_Easy // Podrazumevani nivo ako nešto nije ispravno
	}
}

func ToProtoStatus(status TourStatus) TOURS.TourStatus{
	switch status {
	case Draft:
		return TOURS.TourStatus_Draft
	case Published:
		return TOURS.TourStatus_Published
	case Archived:
		return TOURS.TourStatus_Archived
	default:
		return TOURS.TourStatus_Draft
	}
}

func toProtoTimestamp(t *time.Time) *timestamppb.Timestamp {
    if t == nil {
        return nil
    }
    return timestamppb.New(*t)
}