package api_test

import (
	"testing"

	"github.com/minodisk/qiitactl/api"
)

func TestNewError(t *testing.T) {
	e, err := api.NewError([]byte(`{
		"type": "not_found",
		"message": "Not found"
	}`))
	if err != nil {
		t.Fatal(err)
	}
	if e.Type != "not_found" {
		t.Errorf("wrong type: %s", e.Type)
	}
	if e.Message != "Not found" {
		t.Errorf("wrong message: %s", e.Message)
	}
}

func TestError(t *testing.T) {
	e := api.Error{
		Type:    "not_found",
		Message: "Not found",
	}
	err := e.Error()
	if err.Error() != "not_found: Not found" {
		t.Errorf("wrong Error: %s", err.Error())
	}
}
