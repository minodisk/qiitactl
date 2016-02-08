package model

import "encoding/json"

// Tag is a label for a content of a post.
type Tag struct {
	Name     string   `json:"name"`
	Versions []string `json:"versions"`
}

// MarshalJSON encodes the Tag into JSON formatted bytes.
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
