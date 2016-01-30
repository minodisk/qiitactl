package model

import (
	"testing"
	"time"
)

func TestNewItemWithWrongMeta(t *testing.T) {
	_, err := NewItem([]byte(`XXXXXXXX
<!--
id: abcdefghijklmnopqrst
url: http://example.com/myitem
created_at: 2013-12-10T12:29:14+09:00
updated_at: 2015-02-25T09:26:30+09:00
private: true
coediting: false
tags:
  TypeScript:
  Docker:
    - 1.9
  Go:
    - 1.4.3
    - 1.5.3
-->
# Main title
## Sub title
Paragraph
`))
	if err == nil {
		t.Errorf("start without meta comment should return error")
	}
}

func TestNewItemWithWrongTag(t *testing.T) {
	_, err := NewItem([]byte(`<!--
id: abcdefghijklmnopqrst
url: http://example.com/myitem
created_at: 2013-12-10T12:29:14+09:00
updated_at: 2015-02-25T09:26:30+09:00
private: true
coediting: false
tags:
  TypeScript
  Docker:
    - 1.9
  Go:
    - 1.4.3
    - 1.5.3
-->
# Main title
## Sub title
Paragraph
`))
	if err == nil {
		t.Errorf("should return error with non-object element in tags")
	}
}

func TestNewItemWithWrongTitle(t *testing.T) {
	_, err := NewItem([]byte(`<!--
id: abcdefghijklmnopqrst
url: http://example.com/myitem
created_at: 2013-12-10T12:29:14+09:00
updated_at: 2015-02-25T09:26:30+09:00
private: true
coediting: false
tags:
  TypeScript:
  Docker:
    - 1.9
  Go:
    - 1.4.3
    - 1.5.3
-->
## Sub title
# Main title
Paragraph
`))
	if err == nil {
		t.Errorf("should return error with non-object element in tags")
	}
}

func TestNewItemWithCorrectText(t *testing.T) {
	item, err := NewItem([]byte(`<!--
id: abcdefghijklmnopqrst
url: http://example.com/myitem
created_at: 2013-12-10T12:29:14+09:00
updated_at: 2015-02-25T09:26:30+09:00
private: true
coediting: false
tags:
  TypeScript:
  Docker:
    - 1.9
  Go:
    - 1.4.3
    - 1.5.3
-->
# Main title
## Sub title
Paragraph
`))
	if err != nil {
		t.Fatal(err)
	}

	if item.Meta.Id != "abcdefghijklmnopqrst" {
		t.Errorf("wrong Id")
	}
	if item.Meta.Url != "http://example.com/myitem" {
		t.Errorf("wrong Url")
	}
	if !item.Meta.CreatedAt.Equal(time.Date(2013, 12, 10, 3, 29, 14, 0, time.UTC)) {
		t.Errorf("wrong CreatedAt")
	}
	if !item.Meta.UpdatedAt.Equal(time.Date(2015, 02, 25, 0, 26, 30, 0, time.UTC)) {
		t.Errorf("wrong UpdatedAt")
	}
	if item.Meta.Private != true {
		t.Errorf("wrong Private")
	}
	if item.Meta.Coediting != false {
		t.Errorf("wrong Coediting")
	}
	if len(item.Meta.Tags) != 3 {
		t.Errorf("wrong Tags length: %d", len(item.Meta.Tags))
	} else {
		for _, tag := range item.Meta.Tags {
			switch tag.Name {
			case "TypeScript":
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
	if item.Title != "Main title" {
		t.Errorf("wrong Title")
	}
	if item.Body != "## Sub title\nParagraph" {
		t.Errorf("wrong Body: %s", item.Body)
	}
}
