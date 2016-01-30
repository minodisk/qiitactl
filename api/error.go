package api

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func NewError(b []byte) (e Error, err error) {
	err = json.Unmarshal(b, &e)
	return
}

func (e Error) Error() (err error) {
	err = fmt.Errorf("%s: %s", e.Type, e.Message)
	return
}
