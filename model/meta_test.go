package model_test

import (
	"testing"
	"time"

	"github.com/minodisk/qiitactl/model"
	"github.com/minodisk/qiitactl/testutil"
)

func TestMetaString(t *testing.T) {
	at := model.Time{time.Date(2011, 2, 3, 4, 5, 6, 0, time.UTC)}
	meta := model.Meta{
		ID:        "XXXXXXXX",
		URL:       "http://example.com",
		CreatedAt: at,
		UpdatedAt: at,
		Private:   true,
		Coediting: true,
		Tags: model.Tags{
			model.Tag{
				Name: "Go",
			},
		},
	}
	expected := `id: XXXXXXXX
url: http://example.com
created_at: 2011-02-03T13:05:06+09:00
updated_at: 2011-02-03T13:05:06+09:00
private: true
coediting: true
tags:
- Go`
	actual := meta.String()
	if actual != expected {
		t.Errorf("wrong string:\n%s", testutil.Diff(expected, actual))
	}
}
