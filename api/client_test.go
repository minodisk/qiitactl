package api_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/testutil"
)

var (
	server     *httptest.Server
	rUserAgent = regexp.MustCompile(`qiitactl/\d+\.\d+\.\d+`)
)

func TestMain(m *testing.M) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/echo", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer XXXXXXXXXXXX" {
			b, _ := json.Marshal(api.ResponseError{
				Type:    "unauthorized",
				Message: "Unauthorized",
			})
			w.WriteHeader(401)
			w.Write(b)
			return
		}

		ua := r.Header.Get("User-Agent")
		if !rUserAgent.MatchString(ua) {
			b, _ := json.Marshal(api.ResponseError{
				Type:    "bad_request",
				Message: "Bad Request",
			})
			w.WriteHeader(400)
			w.Write(b)
			return
		}

		defer r.Body.Close()
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			testutil.ResponseError(w, 500, err)
			return
		}

		if len(b) == 0 {
			fmt.Fprintf(w, "%s %s is accepted", r.Method, r.URL)
			return
		}

		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			testutil.ResponseError(w, 400, api.ResponseError{
				Type:    "bad_request",
				Message: "Bad Request",
			})
			return
		}

		var v interface{}
		err = json.Unmarshal(b, &v)
		if err != nil {
			testutil.ResponseError(w, 500, err)
			return
		}
		b, err = json.Marshal(v)
		if err != nil {
			testutil.ResponseError(w, 500, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})

	server = httptest.NewServer(mux)
	defer server.Close()

	code := m.Run()

	// clean up
	testutil.CleanUp()

	os.Exit(code)
}

func TestBuildURL(t *testing.T) {
	func() {
		actual := api.BuildURL("", "/echo")
		if actual != "https://qiita.com/api/v2/echo" {
			t.Errorf("wrong url: %s", actual)
		}
	}()

	func() {
		actual := api.BuildURL("increments", "/echo")
		if actual != "https://increments.qiita.com/api/v2/echo" {
			t.Errorf("wrong url: %s", actual)
		}
	}()
}

func TestNewClientWithEmptyToken(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	_, err := api.NewClient(nil)
	_, ok := err.(api.EmptyTokenError)
	if !ok {
		t.Fatal("empty token error should occur")
	}
}

func TestClientProcess(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client, err := api.NewClient(func(subDomain, path string) (url string) {
		url = fmt.Sprintf("%s%s%s", server.URL, "/api/v2", path)
		return
	})
	if err != nil {
		log.Fatal(err)
	}

	body, _, err := client.Process("OPTIONS", "", "/echo", nil)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != fmt.Sprintf("%s /api/v2%s is accepted", "OPTIONS", "/echo") {
		t.Errorf("wrong body: %s", body)
	}
}

func TestClientProcessWithWrongToken(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	err := os.Setenv("QIITA_ACCESS_TOKEN", "YYYYYYYYYYYY")
	if err != nil {
		log.Fatal(err)
	}
	client, err := api.NewClient(func(subDomain, path string) (url string) {
		url = fmt.Sprintf("%s%s%s", server.URL, "/api/v2", path)
		return
	})
	if err != nil {
		log.Fatal(err)
	}

	_, _, err = client.Process("OPTIONS", "", "/echo", nil)
	_, ok := err.(api.WrongTokenError)
	if !ok {
		t.Fatal("wrong token error should occur")
	}
}

func TestClientProcessWithWrongURL(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client, err := api.NewClient(func(subDomain, path string) (url string) {
		url = fmt.Sprintf("%s%s%s", server.URL, "/api/v2", path)
		return
	})
	if err != nil {
		log.Fatal(err)
	}

	_, _, err = client.Process("GET", "", "/wrong/url", nil)
	_, ok := err.(api.StatusError)
	if !ok {
		t.Fatal("status error should occur")
	}
}

func TestClientPost(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client, err := api.NewClient(func(subDomain, path string) (url string) {
		url = fmt.Sprintf("%s%s%s", server.URL, "/api/v2", path)
		return
	})
	if err != nil {
		log.Fatal(err)
	}

	body, _, err := client.Post("", "/echo", "data")
	if err != nil {
		t.Fatal(err)
	}

	var b string
	err = json.Unmarshal(body, &b)
	if err != nil {
		t.Fatal(err)
	}
	if b != "data" {
		t.Errorf("wrong body: %s", b)
	}
}

func TestClientGet(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client, err := api.NewClient(func(subDomain, path string) (url string) {
		url = fmt.Sprintf("%s%s%s", server.URL, "/api/v2", path)
		return
	})
	if err != nil {
		log.Fatal(err)
	}

	body, _, err := client.Get("", "/echo", nil)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != fmt.Sprintf("%s /api/v2%s is accepted", "GET", "/echo") {
		t.Errorf("wrong body: %s", body)
	}
}

func TestClientPatch(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client, err := api.NewClient(func(subDomain, path string) (url string) {
		url = fmt.Sprintf("%s%s%s", server.URL, "/api/v2", path)
		return
	})
	if err != nil {
		log.Fatal(err)
	}

	body, _, err := client.Patch("", "/echo", nil)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != fmt.Sprintf("%s /api/v2%s is accepted", "PATCH", "/echo") {
		t.Errorf("wrong body: %s", body)
	}
}

func TestClientDelete(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client, err := api.NewClient(func(subDomain, path string) (url string) {
		url = fmt.Sprintf("%s%s%s", server.URL, "/api/v2", path)
		return
	})
	if err != nil {
		log.Fatal(err)
	}

	body, _, err := client.Delete("", "/echo", nil)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != fmt.Sprintf("%s /api/v2%s is accepted", "DELETE", "/echo") {
		t.Errorf("wrong body: %s", body)
	}
}

func TestEmptyTokenError(t *testing.T) {
	err := api.EmptyTokenError{}
	if !strings.HasPrefix(err.Error(), "empty token:") {
		t.Errorf("wrong Error: %s", err.Error())
	}
}

func TestWrongTokenError(t *testing.T) {
	err := api.WrongTokenError{}
	if !strings.HasPrefix(err.Error(), "wrong token:") {
		t.Errorf("wrong Error: %s", err.Error())
	}
}

func TestResponseError(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	err := api.ResponseError{
		Type:    "not_found",
		Message: "Not found",
	}
	if err.Error() != "Not found" {
		t.Errorf("wrong Error: %s", err.Error())
	}
}

func TestStatusError(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	err := api.StatusError{
		Code:    404,
		Message: "404 Not Found",
	}
	if err.Error() != "404 Not Found" {
		t.Errorf("wrong Error: %s", err.Error())
	}
}
