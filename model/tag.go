package model

import "encoding/json"

type Tag struct {
	Name     string   `json:"name"`
	Versions []string `json:"versions"`
}

func (tag Tag) MarshalJSON() ([]byte, error) {
	type T Tag
	t := T{
		Name:     tag.Name,
		Versions: tag.Versions,
	}
	if t.Versions == nil {
		t.Versions = []string{}
	}
	return json.Marshal(t)
}
