package model_test

import (
	"encoding/json"
	"testing"

	"github.com/minodisk/qiitactl/model"
	"github.com/minodisk/qiitactl/testutil"
)

func TestTagMarshalJSONWithNoVersion(t *testing.T) {
	tag := model.Tag{
		Name: "Go",
	}

	b, err := json.MarshalIndent(tag, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	actual := string(b)
	expected := `{
  "name": "Go",
  "versions": []
}`
	if actual != expected {
		t.Errorf("wrong JSON:\n%s", testutil.Diff(expected, actual))
	}
}
