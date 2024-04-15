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
	EncounterID    string				`bson:"encounterId,omitempty" json:"encounterId"`
	CompletionTime string			`bson:"completionTime,omitempty" json:"completionTime"`
	IsCompleted    bool					`bson:"iscompleted" json:"iscompleted"`
}

type EncounterExecutions []*EncounterExecution

func (e *EncounterExecutions) ToJSON(w io.Writer) error {
	d := json.NewEncoder(w)
	return d.Encode(e)
}

func (exec *EncounterExecution) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(exec);
}

func (exec *EncounterExecution) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	if err := d.Decode(exec); err != nil {
		return err
	}

	// Provera da li je CompletionTime nil
	if exec.CompletionTime != "" {
		// Parsiranje stringa u vreme koristeći odgovarajući format
		t, err := time.Parse("2006-01-02T15:04:05", exec.CompletionTime)
		if err != nil {
			return err
		}
		exec.CompletionTime = t.String() // Postavljanje CompletionTime na parsirano vreme
	}

	return nil
}

