package command_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/command"
	"github.com/minodisk/qiitactl/model"
)

var (
	serverMine *httptest.Server
	serverTeam *httptest.Server
	client     api.Client
)

func TestMain(m *testing.M) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/authenticated_user/items", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			log.Fatalf("wrong method: %s", r.Method)
		}
		var body string
		if r.URL.Query().Get("page") == "1" {
			body = `[
				{
					"rendered_body": "<h2>Example body</h2>",
					"body": "## Example body",
					"coediting": false,
					"created_at": "2000-01-01T00:00:00+00:00",
					"id": "4bd431809afb1bb99e4f",
					"private": false,
					"tags": [
						{
							"name": "Ruby",
							"versions": [
								"0.0.1"
							]
						}
					],
					"title": "Example title",
					"updated_at": "2000-01-01T00:00:00+00:00",
					"url": "https://qiita.com/yaotti/items/4bd431809afb1bb99e4f",
					"user": {
						"description": "Hello, world.",
						"facebook_id": "yaotti",
						"followees_count": 100,
						"followers_count": 200,
						"github_login_name": "yaotti",
						"id": "yaotti",
						"items_count": 300,
						"linkedin_id": "yaotti",
						"location": "Tokyo, Japan",
						"name": "Hiroshige Umino",
						"organization": "Increments Inc",
						"permanent_id": 1,
						"profile_image_url": "https://si0.twimg.com/profile_images/2309761038/1ijg13pfs0dg84sk2y0h_normal.jpeg",
						"twitter_screen_name": "yaotti",
						"website_url": "http://yaotti.hatenablog.com"
					}
				}
			]`
		} else {
			body = "[]"
		}
		w.Write([]byte(body))
	})
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
	mux.HandleFunc("/api/v2/items/4bd431809afb1bb99e4f", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL)
		switch r.Method {
		case "PATCH":
			defer r.Body.Close()

			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				responseError(w, 500, err)
				return
			}
			if string(b) == "" {
				responseAPIError(w, 500, api.Error{
					Type:    "fatal",
					Message: "empty body",
				})
				return
			}

			var post model.Post
			err = json.Unmarshal(b, &post)
			if err != nil {
				responseError(w, 500, err)
				return
			}

			post.UpdatedAt = model.Time{Time: time.Date(2016, 2, 1, 12, 51, 42, 0, time.UTC)}
			log.Println(post)
			b, err = json.Marshal(post)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
				return
			}

			_, err = w.Write(b)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
				return
			}
		default:
			w.WriteHeader(405)
		}
	})

	serverMine = httptest.NewServer(mux)
	defer serverMine.Close()

	mux = http.NewServeMux()
	mux.HandleFunc("/api/v2/authenticated_user/items", func(w http.ResponseWriter, r *http.Request) {
		var body string
		if r.URL.Query().Get("page") == "1" {
			body = `[
				{
					"rendered_body": "<h2>Example body in team</h2>",
					"body": "## Example body in team",
					"coediting": false,
					"created_at": "2000-01-01T00:00:00+00:00",
					"id": "4bd431809afb1bb99e4t",
					"private": false,
					"tags": [
						{
							"name": "Ruby",
							"versions": [
								"0.0.1"
							]
						}
					],
					"title": "Example title in team",
					"updated_at": "2000-01-01T00:00:00+00:00",
					"url": "https://qiita.com/yaotti/items/4bd431809afb1bb99e4f",
					"user": {
						"description": "Hello, world.",
						"facebook_id": "yaotti",
						"followees_count": 100,
						"followers_count": 200,
						"github_login_name": "yaotti",
						"id": "yaotti",
						"items_count": 300,
						"linkedin_id": "yaotti",
						"location": "Tokyo, Japan",
						"name": "Hiroshige Umino",
						"organization": "Increments Inc",
						"permanent_id": 1,
						"profile_image_url": "https://si0.twimg.com/profile_images/2309761038/1ijg13pfs0dg84sk2y0h_normal.jpeg",
						"twitter_screen_name": "yaotti",
						"website_url": "http://yaotti.hatenablog.com"
					}
				}
			]`
		} else {
			body = "[]"
		}
		w.Write([]byte(body))
	})
	serverTeam = httptest.NewServer(mux)
	defer serverTeam.Close()

	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}

	client, err = api.NewClient(func(subDomain, path string) (url string) {
		switch subDomain {
		case "":
			url = fmt.Sprintf("%s%s%s", serverMine.URL, "/api/v2", path)
		case "increments":
			url = fmt.Sprintf("%s%s%s", serverTeam.URL, "/api/v2", path)
		default:
			log.Fatalf("wrong sub domain \"%s\"", subDomain)
		}
		return
	})
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	// Clean up trashes
	// os.RemoveAll(model.DirMine)
	// os.RemoveAll("increments")

	os.Exit(code)
}

func responseError(w http.ResponseWriter, statusCode int, err error) {

	responseAPIError(w, statusCode, api.Error{
		Type:    "error",
		Message: err.Error(),
	})
}

func responseAPIError(w http.ResponseWriter, statusCode int, err api.Error) {
	w.WriteHeader(statusCode)
	b, _ := json.Marshal(err)
	w.Write(b)
}

func TestShowPosts(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	err := command.ShowPosts(client, buf)
	if err != nil {
		t.Fatal(err)
	}

	if string(buf.Bytes()) != `Posts in Qiita:
4bd431809afb1bb99e4f 2000/01/01 Example title
Posts in Qiita:Team (Increments Inc.):
4bd431809afb1bb99e4t 2000/01/01 Example title in team
` {
		t.Errorf("written text is wrong: %s", buf.Bytes())
	}
}

func TestFetchPosts(t *testing.T) {
	err := command.FetchPosts(client)
	if err != nil {
		t.Fatal(err)
	}

	filepath.Walk(model.DirMine, func(path string, info os.FileInfo, e error) (err error) {
		if e != nil {
			t.Fatal(e)
		}
		if !info.IsDir() {
			if path != fmt.Sprintf("%s/2000-01-01-example-title.md", model.DirMine) {
				t.Fatalf("wrong file is created \"%s\"", path)
			}
			b, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if string(b) != `<!--
id: 4bd431809afb1bb99e4f
url: https://qiita.com/yaotti/items/4bd431809afb1bb99e4f
created_at: 2000-01-01T09:00:00+09:00
updated_at: 2000-01-01T09:00:00+09:00
private: false
coediting: false
tags:
  Ruby:
  - 0.0.1
-->
# Example title
## Example body` {
				t.Errorf("wrong body \"%s\"", b)
			}
		}
		return
	})

	filepath.Walk("increments", func(path string, info os.FileInfo, e error) (err error) {
		if e != nil {
			t.Fatal(e)
		}
		if !info.IsDir() {
			if path != "increments/2000-01-01-example-title-in-team.md" {
				t.Fatalf("wrong file is created \"%s\"", path)
			}
			b, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if string(b) != `<!--
id: 4bd431809afb1bb99e4t
url: https://qiita.com/yaotti/items/4bd431809afb1bb99e4f
created_at: 2000-01-01T09:00:00+09:00
updated_at: 2000-01-01T09:00:00+09:00
private: false
coediting: false
tags:
  Ruby:
  - 0.0.1
-->
# Example title in team
## Example body in team` {
				t.Errorf("wrong body \"%s\"", b)
			}
		}
		return
	})
}

func TestUpdatePost(t *testing.T) {
	path := fmt.Sprintf("%s/2000-01-01-example-title.md", model.DirMine)
	err := ioutil.WriteFile(path, []byte(`<!--
id: 4bd431809afb1bb99e4f
url: https://qiita.com/yaotti/items/4bd431809afb1bb99e4f
created_at: 2000-01-01T09:00:00+09:00
updated_at: 2000-01-01T09:00:00+09:00
private: false
coediting: false
tags:
  Ruby:
  - 0.0.1
  NewTag:
  - 1.0
-->
# Example new title
## Example new body`), 664)
	if err != nil {
		t.Fatal(err)
	}

	err = command.UpdatePost(client, path)
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if string(b) != `<!--
id: 4bd431809afb1bb99e4f
url: https://qiita.com/yaotti/items/4bd431809afb1bb99e4f
created_at: 2000-01-01T09:00:00+09:00
updated_at: 2016-02-01T21:51:42+09:00
private: false
coediting: false
tags:
  NewTag:
  - "1.0"
  Ruby:
  - 0.0.1
-->
# Example new title
## Example new body` {
		t.Errorf("wrong content: %s", b)
	}
}
