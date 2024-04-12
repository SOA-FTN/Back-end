package model

import (
	"encoding/json"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EncounterExecution struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserID         int64				`bson:"userId,omitempty" json:"userId"`
	EncounterID    int64				`bson:"encounterId,omitempty" json:"encounterId"`
	CompletionTime *time.Time			`bson:"completionTime,omitempty" json:"completionTime"`
	IsCompleted    bool					`bson:"completed,omitempty" json:"completed"`
}


func (exec *EncounterExecution) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(exec);
}

func (exec *EncounterExecution) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(exec)
}
