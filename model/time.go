package model

import (
	"encoding/json"
	"time"
)

// Time is a wrapper of time.Time.
type Time struct {
	time.Time
}

// ITime should have a format function.
type ITime interface {
	Format(layout string) string
}

// UnmarshalJSON unmarshals time as JSON.
func (t *Time) UnmarshalJSON(b []byte) (err error) {
	var s string
	err = json.Unmarshal(b, &s)
	if err != nil {
		return
	}
	t.Time, err = time.Parse(time.RFC3339, s)
	if err != nil {
		return
	}
	return
}

// MarshalYAML marshals time as YAML.
func (t Time) MarshalYAML() (data interface{}, err error) {
	data = t.Local().Format(time.RFC3339)
	return
}

// FormatDate formats only date as slashed style.
func (t Time) FormatDate() (s string) {
	s = t.Local().Format("2006/01/02")
	return
}
