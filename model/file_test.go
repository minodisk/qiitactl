package model_test

import (
	"testing"
	"time"

	"github.com/minodisk/qiitactl/model"
)

func TestFillPath(t *testing.T) {
	var f model.File
	f.FillPath(time.Date(2016, 2, 3, 4, 5, 6, 0, time.UTC), "", nil)
	if f.Path != "mine/2016/02/03-.md" {
		t.Errorf("wrong path: %s", f.Path)
	}
}
