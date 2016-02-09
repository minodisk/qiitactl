package model_test

import (
	"testing"
	"time"

	"github.com/minodisk/qiitactl/model"
	"github.com/minodisk/qiitactl/testutil"
)

func TestMetaEncode(t *testing.T) {
	at := model.Time{Time: time.Date(2011, 2, 3, 4, 5, 6, 0, time.UTC)}
	tags := model.Tags{
		model.Tag{
			Name: "Go",
		},
	}
	meta := model.Meta{
		ID:        "4bd431809afb1bb99e4f",
		URL:       "https://qiita.com/yaotti/items/4bd431809afb1bb99e4f",
		CreatedAt: at,
		UpdatedAt: at,
		Private:   true,
		Coediting: true,
		Tags:      tags,
		Team: &model.Team{
			Active: true,
			ID:     "increments",
			Name:   "Increments Inc.",
		},
	}
	expected := `id: 4bd431809afb1bb99e4f
url: https://qiita.com/yaotti/items/4bd431809afb1bb99e4f
created_at: 2011-02-03T13:05:06+09:00
updated_at: 2011-02-03T13:05:06+09:00
private: true
coediting: true
tags:
- Go
team:
  active: true
  id: increments
  name: Increments Inc.`
	actual := meta.Encode()
	if actual != expected {
		t.Errorf("wrong string:\n%s", testutil.Diff(expected, actual))
	}
}

func TestMetaDecode(t *testing.T) {
	var meta model.Meta
	meta.Decode(`id: 4bd431809afb1bb99e4f
url: https://qiita.com/yaotti/items/4bd431809afb1bb99e4f
created_at: 2013-12-10T12:29:14+09:00
updated_at: 2015-02-25T09:26:30+09:00
private: true
coediting: false
tags:
- JavaScript
- Docker:
  - 1.9
- Go:
  - 1.4.3
  - 1.5.3
team:
  active: true
  id: increments
  name: Increments Inc.`)
	if meta.ID != "4bd431809afb1bb99e4f" {
		t.Errorf("wrong Id")
	}
	if meta.URL != "https://qiita.com/yaotti/items/4bd431809afb1bb99e4f" {
		t.Errorf("wrong Url")
	}
	if !meta.CreatedAt.Equal(time.Date(2013, 12, 10, 3, 29, 14, 0, time.UTC)) {
		t.Errorf("wrong CreatedAt")
	}
	if !meta.UpdatedAt.Equal(time.Date(2015, 02, 25, 0, 26, 30, 0, time.UTC)) {
		t.Errorf("wrong UpdatedAt")
	}
	if meta.Team.Active != true || meta.Team.ID != "increments" || meta.Team.Name != "Increments Inc." {
		t.Errorf("wrong Team")
	}
	if meta.Private != true {
		t.Errorf("wrong Private")
	}
	if meta.Coediting != false {
		t.Errorf("wrong Coediting")
	}
	if len(meta.Tags) != 3 {
		t.Errorf("wrong Tags length: %d", len(meta.Tags))
	} else {
		for _, tag := range meta.Tags {
			switch tag.Name {
			case "JavaScript":
				if len(tag.Versions) != 0 {
					t.Errorf("wrong Tag with no version: %+v", tag)
				}
			case "Docker":
				if len(tag.Versions) != 1 || tag.Versions[0] != "1.9" {
					t.Errorf("wrong Tag with single version: %+v", tag)
				}
			case "Go":
				if len(tag.Versions) != 2 || tag.Versions[0] != "1.4.3" || tag.Versions[1] != "1.5.3" {
					t.Errorf("wrong Tag with multi versions: %+v", tag)
				}
			}
		}
	}
}
