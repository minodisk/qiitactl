package model

import (
	"encoding/json"
	"time"
)

// 出力時に必要に応じてローカルタイムに変換する
type Time struct {
	time.Time
}

type ITime interface {
	Format(layout string) string
}

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

func (t Time) MarshalYAML() (data interface{}, err error) {
	data = t.Local().Format(time.RFC3339)
	return
}

func (t Time) FormatDate() (s string) {
	s = t.Local().Format("2006/01/02")
	return
}
