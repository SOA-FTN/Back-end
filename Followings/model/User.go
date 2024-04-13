package model

import (
	"encoding/json"
	"io"
)

type User struct {
	UserName string `json:"Username"`
}

func (o *User) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(o)
}
