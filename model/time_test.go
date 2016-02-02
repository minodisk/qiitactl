package model_test

import (
	"encoding/json"
	"testing"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/minodisk/qiitactl/model"
)

func TestTime_UnmarshalJSON(t *testing.T) {
	s := `"2016-02-03T13:05:06+09:00"`
	var at model.Time
	err := json.Unmarshal([]byte(s), &at)
	if err != nil {
		t.Fatal(err)
	}
	if at.Year() != 2016 || at.Month() != 2 || at.Day() != 3 || at.Hour() != 13 || at.Minute() != 5 || at.Second() != 6 {
		t.Errorf("wrong unmarshaled json: %s", at)
	}
}

func TestTime_MarshalYAML(t *testing.T) {
	at := model.Time{time.Date(2016, 2, 3, 4, 5, 6, 0, time.UTC)}
	b, err := yaml.Marshal(at)
	if err != nil {
		t.Fatal(err)
	}
	actual := string(b)
	expected := `2016-02-03T13:05:06+09:00
`
	if actual != expected {
		t.Errorf("wrong marshaled yaml: expected %s, but actual %s", expected, actual)
	}
}

func TestTime_FormatDate(t *testing.T) {
	at := model.Time{time.Date(2016, 2, 3, 4, 5, 6, 0, time.UTC)}
	actual := at.FormatDate()
	expected := "2016/02/03"
	if actual != expected {
		t.Errorf("wrong result: expected %s, but actual %s", expected, actual)
	}
}
