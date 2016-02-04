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
	"testing"
	"time"

	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/command"
	"github.com/minodisk/qiitactl/model"
	"github.com/minodisk/qiitactl/testutil"
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
		switch r.Method {
		case "GET":
			w.Write([]byte(`{
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
				}`))

		case "PATCH":
			defer r.Body.Close()

			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				responseError(w, 500, err)
				return
			}
			if string(b) == "" {
				responseAPIError(w, 500, api.ResponseError{
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
					"created_at": "2015-09-25T00:00:00+00:00",
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
					"updated_at": "2015-09-25T00:00:00+00:00",
					"url": "https://increments.qiita.com/yaotti/items/4bd431809afb1bb99e4t",
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
	cleanUp()

	os.Exit(code)
}

func responseError(w http.ResponseWriter, statusCode int, err error) {

	responseAPIError(w, statusCode, api.ResponseError{
		Type:    "error",
		Message: err.Error(),
	})
}

func responseAPIError(w http.ResponseWriter, statusCode int, err api.ResponseError) {
	w.WriteHeader(statusCode)
	b, _ := json.Marshal(err)
	w.Write(b)
}

func TestFetchPostWithID(t *testing.T) {
	defer cleanUp()

	err := command.FetchPost(client, "4bd431809afb1bb99e4f", "")
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadFile("mine/2000/01/01-example-title.md")
	if err != nil {
		t.Fatal(err)
	}
	actual := string(b)
	expected := `<!--
id: 4bd431809afb1bb99e4f
url: https://qiita.com/yaotti/items/4bd431809afb1bb99e4f
created_at: 2000-01-01T09:00:00+09:00
updated_at: 2000-01-01T09:00:00+09:00
private: false
coediting: false
tags:
- Ruby:
  - 0.0.1
-->
# Example title
## Example body`
	if actual != expected {
		t.Errorf("wrong body:\n%s", testutil.Diff(expected, actual))
	}
}

func TestFetchPostWithFilename(t *testing.T) {
	defer cleanUp()

	err := command.FetchPost(client, "4bd431809afb1bb99e4f", "")
	if err != nil {
		t.Fatal(err)
	}

	err = command.FetchPost(client, "", "mine/2000/01/01-example-title.md")
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadFile("mine/2000/01/01-example-title.md")
	if err != nil {
		t.Fatal(err)
	}
	actual := string(b)
	expected := `<!--
id: 4bd431809afb1bb99e4f
url: https://qiita.com/yaotti/items/4bd431809afb1bb99e4f
created_at: 2000-01-01T09:00:00+09:00
updated_at: 2000-01-01T09:00:00+09:00
private: false
coediting: false
tags:
- Ruby:
  - 0.0.1
-->
# Example title
## Example body`
	if actual != expected {
		t.Errorf("wrong body:\n%s", testutil.Diff(expected, actual))
	}
}

func TestShowPostWithID(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	err := command.ShowPost(client, buf, "4bd431809afb1bb99e4f", "")
	if err != nil {
		t.Fatal(err)
	}

	if string(buf.Bytes()) != `4bd431809afb1bb99e4f 2000/01/01 Example title
` {
		t.Errorf("written text is wrong: %s", buf.Bytes())
	}
}

func TestShowPostWithFilename(t *testing.T) {
	defer cleanUp()

	err := command.FetchPost(client, "4bd431809afb1bb99e4f", "")
	if err != nil {
		t.Fatal(err)
	}

	buf := bytes.NewBuffer([]byte{})
	err = command.ShowPost(client, buf, "", "mine/2000/01/01-example-title.md")
	if err != nil {
		t.Fatal(err)
	}

	if string(buf.Bytes()) != `4bd431809afb1bb99e4f 2000/01/01 Example title
` {
		t.Errorf("written text is wrong: %s", buf.Bytes())
	}
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
4bd431809afb1bb99e4t 2015/09/25 Example title in team
` {
		t.Errorf("written text is wrong: %s", buf.Bytes())
	}
}

func TestFetchPosts(t *testing.T) {
	defer cleanUp()

	err := command.FetchPosts(client)
	if err != nil {
		t.Fatal(err)
	}

	func() {
		b, err := ioutil.ReadFile("mine/2000/01/01-example-title.md")
		if err != nil {
			t.Fatal(err)
		}
		actual := string(b)
		expected := `<!--
id: 4bd431809afb1bb99e4f
url: https://qiita.com/yaotti/items/4bd431809afb1bb99e4f
created_at: 2000-01-01T09:00:00+09:00
updated_at: 2000-01-01T09:00:00+09:00
private: false
coediting: false
tags:
- Ruby:
  - 0.0.1
-->
# Example title
## Example body`
		if actual != expected {
			t.Errorf("wrong body:\n%s", testutil.Diff(expected, actual))
		}
	}()

	func() {
		b, err := ioutil.ReadFile("increments/2015/09/25-example-title-in-team.md")
		if err != nil {
			t.Fatal(err)
		}
		actual := string(b)
		expected := `<!--
id: 4bd431809afb1bb99e4t
url: https://increments.qiita.com/yaotti/items/4bd431809afb1bb99e4t
created_at: 2015-09-25T09:00:00+09:00
updated_at: 2015-09-25T09:00:00+09:00
private: false
coediting: false
tags:
- Ruby:
  - 0.0.1
-->
# Example title in team
## Example body in team`
		if actual != expected {
			t.Errorf("wrong body:\n%s", testutil.Diff(expected, actual))
		}
	}()
}

func TestUpdatePost(t *testing.T) {
	defer cleanUp()

	path := fmt.Sprintf("%s/2000/01/01-example-title.md", model.DirMine)

	err := command.FetchPost(client, "4bd431809afb1bb99e4f", "")
	if err != nil {
		t.Fatal(err)
	}

	err = ioutil.WriteFile(path, []byte(`<!--
id: 4bd431809afb1bb99e4f
url: https://qiita.com/yaotti/items/4bd431809afb1bb99e4f
created_at: 2000-01-01T09:00:00+09:00
updated_at: 2000-01-01T09:00:00+09:00
private: false
coediting: false
tags:
- Ruby:
  - 0.0.1
- NewTag:
  - "1.0"
-->
# Example edited title
## Example edited body`), 0664)
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

	actual := string(b)
	expected := `<!--
id: 4bd431809afb1bb99e4f
url: https://qiita.com/yaotti/items/4bd431809afb1bb99e4f
created_at: 2000-01-01T09:00:00+09:00
updated_at: 2016-02-01T21:51:42+09:00
private: false
coediting: false
tags:
- Ruby:
  - 0.0.1
- NewTag:
  - "1.0"
-->
# Example edited title
## Example edited body`
	if actual != expected {
		t.Errorf("wrong content:\n%s", testutil.Diff(expected, actual))
	}
}

// func TestTracingPostFileWithPostID(t *testing.T) {
// 	defer cleanUp()
//
// 	err := os.MkdirAll("mine/1990/07", 0755)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	err = ioutil.WriteFile("mine/1990/07/03-example-old-title.md", []byte(`<!--
// id: 4bd431809afb1bb99e4f
// url: https://qiita.com/yaotti/items/4bd431809afb1bb99e4f
// created_at: 1990-07-03T09:00:00+09:00
// updated_at: 1990-07-03T09:00:00+09:00
// private: false
// coediting: false
// tags:
// - Ruby:
//   - 0.0.0
// -->
// # Example old title
// ## Example old body`), 0664)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	err = command.FetchPost(client, "4bd431809afb1bb99e4f", "")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	b, err := ioutil.ReadFile("mine/1990/07/03-example-old-title.md")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	actual := string(b)
// 	expected := `<!--
// id: 4bd431809afb1bb99e4f
// url: https://qiita.com/yaotti/items/4bd431809afb1bb99e4f
// created_at: 2000-01-01T09:00:00+09:00
// updated_at: 2000-01-01T09:00:00+09:00
// private: false
// coediting: false
// tags:
// - Ruby:
//   - 0.0.1
// -->
// # Example title
// ## Example body`
// 	if actual != expected {
// 		t.Errorf("wrong body:\n%s", testutil.Diff(expected, actual))
// 	}
// }
//
// func TestTracingPostFilesWithPostIDs(t *testing.T) {
// 	defer cleanUp()
//
// 	err := os.MkdirAll("mine/1990/07", 0755)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	err = ioutil.WriteFile("mine/1990/07/03-example-old-title.md", []byte(`<!--
// id: 4bd431809afb1bb99e4f
// url: https://qiita.com/yaotti/items/4bd431809afb1bb99e4f
// created_at: 1990-07-03T09:00:00+09:00
// updated_at: 1990-07-03T09:00:00+09:00
// private: false
// coediting: false
// tags:
// - Ruby:
//   - 0.0.0
// -->
// # Example old title
// ## Example old body`), 0664)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	err = os.MkdirAll("increments/1993/12", 0755)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	err = ioutil.WriteFile("increments/1993/12/30-example-old-title-in-team.md", []byte(`<!--
// id: 4bd431809afb1bb99e4t
// url: https://increments.qiita.com/yaotti/items/4bd431809afb1bb99e4t
// created_at: 1993-12-30:22:11:31+09:00
// updated_at: 1993-12-30:22:11:31+09:00
// private: false
// coediting: false
// tags:
// - Ruby:
//   - 0.0.1
// -->
// # Example old title in team
// ## Example old body in team`), 0664)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	err = command.FetchPosts(client)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	b, err := ioutil.ReadFile("mine/1990/07/03-example-old-title.md")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	actual := string(b)
// 	expected := `<!--
// id: 4bd431809afb1bb99e4f
// url: https://qiita.com/yaotti/items/4bd431809afb1bb99e4f
// created_at: 2000-01-01T09:00:00+09:00
// updated_at: 2000-01-01T09:00:00+09:00
// private: false
// coediting: false
// tags:
// - Ruby:
//   - 0.0.1
// -->
// # Example title
// ## Example body`
// 	if actual != expected {
// 		t.Errorf("wrong body:\n%s", testutil.Diff(expected, actual))
// 	}
//
// 	b, err = ioutil.ReadFile("increments/1993/12/30-example-old-title-in-team.md")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	actual = string(b)
// 	expected = `<!--
// id: 4bd431809afb1bb99e4t
// url: https://increments.qiita.com/yaotti/items/4bd431809afb1bb99e4t
// created_at: 2015-09-25T09:00:00+09:00
// updated_at: 2015-09-25T09:00:00+09:00
// private: false
// coediting: false
// tags:
// - Ruby:
//   - 0.0.1
// -->
// # Example title in team
// ## Example body in team`
// 	if actual != expected {
// 		t.Errorf("wrong body:\n%s", testutil.Diff(expected, actual))
// 	}
// }

func cleanUp() {
	os.RemoveAll("mine")
	os.RemoveAll("increments")
}
