package model_test

import (
	"testing"
	"time"

	"github.com/minodisk/qiitactl/model"
)

func TestFillPath_NoTitle(t *testing.T) {
	var f model.File
	f.FillPath(time.Date(2016, 2, 3, 4, 5, 6, 0, time.UTC), "", nil)
	if f.Path != "mine/2016/02/03.md" {
		t.Errorf("wrong path: %s", f.Path)
	}
}

func TestFillPath_NormalTitle(t *testing.T) {
	var f model.File
	f.FillPath(time.Date(2016, 2, 3, 4, 5, 6, 0, time.UTC), "title", nil)
	if f.Path != "mine/2016/02/03-title.md" {
		t.Errorf("wrong path: %s", f.Path)
	}
}

func TestFillPath_TitleWithSpace(t *testing.T) {
	var f model.File
	f.FillPath(time.Date(2016, 2, 3, 4, 5, 6, 0, time.UTC), "example title", nil)
	if f.Path != "mine/2016/02/03-example-title.md" {
		t.Errorf("wrong path: %s", f.Path)
	}
}

func TestFillPath_TitleWithUnicode(t *testing.T) {
	var f model.File
	f.FillPath(time.Date(2016, 2, 3, 4, 5, 6, 0, time.UTC), "example タイトル", nil)
	if f.Path != "mine/2016/02/03-example.md" {
		t.Errorf("wrong path: %s", f.Path)
	}
	f.FillPath(time.Date(2016, 2, 3, 4, 5, 6, 0, time.UTC), "例 title", nil)
	if f.Path != "mine/2016/02/03-title.md" {
		t.Errorf("wrong path: %s", f.Path)
	}
}

func TestFillPath_TitleWithInvalidCharacter(t *testing.T) {
	var f model.File
	f.FillPath(time.Date(2016, 2, 3, 4, 5, 6, 0, time.UTC), "example, title", nil)
	if f.Path != "mine/2016/02/03-example-title.md" {
		t.Errorf("wrong path: %s", f.Path)
	}
	f.FillPath(time.Date(2016, 2, 3, 4, 5, 6, 0, time.UTC), "example::title", nil)
	if f.Path != "mine/2016/02/03-example-title.md" {
		t.Errorf("wrong path: %s", f.Path)
	}
}
