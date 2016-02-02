package api_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"github.com/minodisk/qiitactl/api"
)

var (
	server     *httptest.Server
	client     api.Client
	rUserAgent = regexp.MustCompile(`qiitactl/\d+\.\d+\.\d+`)
)

func TestMain(m *testing.M) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer XXXXXXXXXXXX" {
			b, _ := json.Marshal(api.Error{
				Type:    "unauthorized",
				Message: "Unauthorized",
			})
			w.WriteHeader(401)
			w.Write(b)
			return
		}
		ua := r.Header.Get("User-Agent")
		if !rUserAgent.MatchString(ua) {
			b, _ := json.Marshal(api.Error{
				Type:    "bad_request",
				Message: "Bad Request",
			})
			w.WriteHeader(400)
			w.Write(b)
			return
		}
	})

	server = httptest.NewServer(mux)
	defer server.Close()

	var err error

	err = os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}

	client, err = api.NewClient(func(subDomain, path string) (url string) {
		url = fmt.Sprintf("%s%s%s", server.URL, "/api/v2", path)
		return
	})
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	// clean up

	os.Exit(code)
}

func TestProcess(t *testing.T) {
	_, err := client.Process("OPTIONS", "", "", nil)
	if err != nil {
		t.Fatal(err)
	}
}
