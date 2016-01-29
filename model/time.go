package model

import (
	"encoding/json"
	"time"
)

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) (err error) {
	var s string
	err = json.Unmarshal(b, &s)
	if err != nil {
		return
	}
	ti, err := time.Parse(time.RFC3339, s)
	t.Time = ti.Local()
	if err != nil {
		return
	}
	return
}

func (t *Time) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var s string
	err = unmarshal(&s)
	if err != nil {
		return
	}
	ti, err := time.Parse(time.RFC3339, s)
	t.Time = ti.Local()
	if err != nil {
		return
	}
	return
}

func (t Time) Format() (s string) {
	s = t.Time.Format(time.RFC3339)
	return
}

func (t Time) FormatDate() (s string) {
	s = t.Time.Format("2006/01/02")
	return
}

func (t Time) ToPath() (path string) {
	path = t.Time.Format("2006-01-02")
	return
}
