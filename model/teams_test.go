package model_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/model"
)

func TestTeams_FetchTeams(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/teams", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			log.Fatalf("wrong method: %s", r.Method)
		}
		w.Write([]byte(`[
			{
				"active": true,
				"id": "increments",
				"name": "Increments Inc."
			}
		]`))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		t.Fatal(err)
	}

	client, err := api.NewClient(func(subDomain, path string) (url string) {
		url = fmt.Sprintf("%s%s%s", server.URL, "/api/v2", path)
		return
	})
	if err != nil {
		t.Fatal(err)
	}

	teams, err := model.FetchTeams(client)
	if err != nil {
		t.Fatal(err)
	}

	if len(teams) != 1 {
		t.Fatalf("wrong length of teams: expected %d, but actual %d", 1, len(teams))
	}

	if teams[0].Active != true || teams[0].ID != "increments" || teams[0].Name != "Increments Inc." {
		t.Errorf("wrong team: %+v", teams[0])
	}
}

func TestTeams_FetchTeamsError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/teams", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(503)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		t.Fatal(err)
	}

	client, err := api.NewClient(func(subDomain, path string) (url string) {
		url = fmt.Sprintf("%s%s%s", server.URL, "/api/v2", path)
		return
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = model.FetchTeams(client)
	if err == nil {
		t.Fatal("should occur error")
	}
}
