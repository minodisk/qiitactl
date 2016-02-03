package model_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/model"
	"github.com/minodisk/qiitactl/testutil"
)

func TestFetchPost(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/items/4bd431809afb1bb99e4f", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(405)
			b, _ := json.Marshal(api.ResponseError{"method_not_allowed", "Method Not Allowed"})
			w.Write(b)
			return
		}

		body := `{
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
		}`
		w.Write([]byte(body))
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

	team := model.Team{
		Active: true,
		ID:     "increments",
		Name:   "Increments Inc",
	}
	post, err := model.FetchPost(client, &team, "4bd431809afb1bb99e4f")
	if err != nil {
		t.Fatal(err)
	}

	if post.RenderedBody != "<h2>Example body</h2>" {
		t.Errorf("wrong RenderedBody: %s", post.RenderedBody)
	}
	if post.Body != "## Example body" {
		t.Errorf("wrong Body: %s", post.Body)
	}
	if post.Coediting != false {
		t.Errorf("wrong Coediting: %b", post.Coediting)
	}
	if !post.CreatedAt.Equal(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("wrong CreatedAt: %s", post.CreatedAt)
	}
	if post.ID != "4bd431809afb1bb99e4f" {
		t.Errorf("wrong ID: %s", post.ID)
	}
	if post.Private != false {
		t.Errorf("wrong Private: %b", post.Private)
	}
	if post.Title != "Example title" {
		t.Errorf("wrong Title: %s", post.Title)
	}
	if !post.UpdatedAt.Equal(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("wrong UpdatedAt: %s", post.UpdatedAt)
	}
	if post.URL != "https://qiita.com/yaotti/items/4bd431809afb1bb99e4f" {
		t.Errorf("wrong URL: %s", post.URL)
	}
	if post.User.Description != "Hello, world." {
		t.Errorf("wrong Description: %s", post.User.Description)
	}
	if post.User.FacebookId != "yaotti" {
		t.Errorf("wrong FacebookId: %s", post.User.FacebookId)
	}
	if post.User.FolloweesCount != 100 {
		t.Errorf("wrong FolloweesCount: %s", post.User.FolloweesCount)
	}
	if post.User.FollowersCount != 200 {
		t.Errorf("wrong FollowersCount: %s", post.User.FollowersCount)
	}
	if post.User.GithubLoginName != "yaotti" {
		t.Errorf("wrong GithubLoginName: %s", post.User.GithubLoginName)
	}
	if post.User.Id != "yaotti" {
		t.Errorf("wrong Id: %s", post.User.Id)
	}
	if post.User.ItemsCount != 300 {
		t.Errorf("wrong ItemsCount: %d", post.User.ItemsCount)
	}
	if post.User.LinkedinId != "yaotti" {
		t.Errorf("wrong LinkedinId: %s", post.User.LinkedinId)
	}
	if post.User.Location != "Tokyo, Japan" {
		t.Errorf("wrong Location: %s", post.User.Location)
	}
	if post.User.Name != "Hiroshige Umino" {
		t.Errorf("wrong Name: %s", post.User.Name)
	}
	if post.User.Organization != "Increments Inc" {
		t.Errorf("wrong Organization: %s", post.User.Organization)
	}
	if post.User.PermanentId != 1 {
		t.Errorf("wrong PermanentId: %d", post.User.PermanentId)
	}
	if post.User.ProfileImageUrl != "https://si0.twimg.com/profile_images/2309761038/1ijg13pfs0dg84sk2y0h_normal.jpeg" {
		t.Errorf("wrong ProfileImageUrl: %s", post.User.ProfileImageUrl)
	}
	if post.User.TwitterScreenName != "yaotti" {
		t.Errorf("wrong TwitterScreenName: %s", post.User.TwitterScreenName)
	}
	if post.User.WebsiteUrl != "http://yaotti.hatenablog.com" {
		t.Errorf("wrong WebsiteUrl: %s", post.User.WebsiteUrl)
	}
	if len(post.Tags) != 1 {
		t.Fatalf("wrong Tags length: %d", len(post.Tags))
	}
	if post.Tags[0].Name != "Ruby" {
		t.Errorf("wrong tag Name: %s", post.Tags[0].Name)
	}
	if len(post.Tags[0].Versions) != 1 {
		t.Fatalf("wrong tag Versions length: %d", len(post.Tags[0].Versions))
	}
	if post.Tags[0].Versions[0] != "0.0.1" {
		t.Errorf("wrong tag Versions: %s", post.Tags[0].Versions[0])
	}
}

func TestFetchPost_ResponseError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/items/4bd431809afb1bb99e4f", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte(`{
  "message": "Not found",
  "type": "not_found"
}`))
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

	_, err = model.FetchPost(client, nil, "4bd431809afb1bb99e4f")
	if err == nil {
		t.Fatal("error should occur")
	}
	_, ok := err.(api.ResponseError)
	if !ok {
		t.Fatalf("wrong type error: %s", reflect.TypeOf(err))
	}
}

func TestFetchPost_StatusError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v2/items/4bd431809afb1bb99e4f", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
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

	_, err = model.FetchPost(client, nil, "4bd431809afb1bb99e4f")
	if err == nil {
		t.Fatal("error should occur")
	}
	_, ok := err.(api.StatusError)
	if !ok {
		t.Fatalf("wrong type error: %s", reflect.TypeOf(err))
	}
}

func TestPostSave(t *testing.T) {
	post := model.NewPost("Example Title 0", &model.Time{time.Date(2015, 11, 28, 13, 2, 37, 0, time.UTC)}, nil)
	err := post.Save()
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll("mine")

	func() {
		a, err := ioutil.ReadFile("mine/2015/11/28-example-title-0.md")
		if err != nil {
			t.Fatal(err)
		}
		actual := string(a)
		expected := `<!--
id: ""
url: ""
created_at: 2015-11-28T22:02:37+09:00
updated_at: 2015-11-28T22:02:37+09:00
private: false
coediting: false
tags: []
-->
# Example Title 0
`
		if actual != expected {
			t.Errorf("wrong content:\n%s", testutil.Diff(expected, actual))
		}
	}()
}

func TestPostEncodeWithNewPost(t *testing.T) {
	post := model.NewPost("Example title", &model.Time{time.Date(2016, 2, 2, 6, 30, 46, 0, time.UTC)}, nil)
	buf := bytes.NewBuffer([]byte{})
	err := post.Encode(buf)
	if err != nil {
		t.Fatal(err)
	}
	actual := string(buf.Bytes())
	expected := `<!--
id: ""
url: ""
created_at: 2016-02-02T15:30:46+09:00
updated_at: 2016-02-02T15:30:46+09:00
private: false
coediting: false
tags: []
-->
# Example title
`
	if expected != actual {
		t.Errorf("wrong content:\n%s", testutil.Diff(expected, actual))
	}
}

func TestPostDecodeWithWrongMeta(t *testing.T) {
	var post model.Post
	err := post.Decode([]byte(`XXXXXXXX
<!--
id: abcdefghijklmnopqrst
url: http://example.com/mypost
created_at: 2013-12-10T12:29:14+09:00
updated_at: 2015-02-25T09:26:30+09:00
private: true
coediting: false
tags:
- TypeScript
- Docker:
  - 1.9
- Go:
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

func TestPostDecodeWithWrongTag(t *testing.T) {
	var post model.Post
	err := post.Decode([]byte(`<!--
id: abcdefghijklmnopqrst
url: http://example.com/mypost
created_at: 2013-12-10T12:29:14+09:00
updated_at: 2015-02-25T09:26:30+09:00
private: true
coediting: false
tags:
  TypeScript: []
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
		t.Errorf("should return error objective tags")
	}
}

func TestDecodeWithWrongTitle(t *testing.T) {
	var post model.Post
	err := post.Decode([]byte(`<!--
id: abcdefghijklmnopqrst
url: http://example.com/mypost
created_at: 2013-12-10T12:29:14+09:00
updated_at: 2015-02-25T09:26:30+09:00
private: true
coediting: false
tags:
- TypeScript
- Docker:
  - 1.9
- Go:
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

func TestPostDecodeWithCorrectMarkdown(t *testing.T) {
	var post model.Post
	err := post.Decode([]byte(`<!--
id: abcdefghijklmnopqrst
url: http://example.com/mypost
created_at: 2013-12-10T12:29:14+09:00
updated_at: 2015-02-25T09:26:30+09:00
private: true
coediting: false
tags:
- TypeScript
- Docker:
  - 1.9
- Go:
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

	if post.Meta.ID != "abcdefghijklmnopqrst" {
		t.Errorf("wrong Id")
	}
	if post.Meta.URL != "http://example.com/mypost" {
		t.Errorf("wrong Url")
	}
	if !post.Meta.CreatedAt.Equal(time.Date(2013, 12, 10, 3, 29, 14, 0, time.UTC)) {
		t.Errorf("wrong CreatedAt")
	}
	if !post.Meta.UpdatedAt.Equal(time.Date(2015, 02, 25, 0, 26, 30, 0, time.UTC)) {
		t.Errorf("wrong UpdatedAt")
	}
	if post.Meta.Private != true {
		t.Errorf("wrong Private")
	}
	if post.Meta.Coediting != false {
		t.Errorf("wrong Coediting")
	}
	if len(post.Meta.Tags) != 3 {
		t.Errorf("wrong Tags length: %d", len(post.Meta.Tags))
	} else {
		for _, tag := range post.Meta.Tags {
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
	if post.Title != "Main title" {
		t.Errorf("wrong Title")
	}
	if post.Body != "## Sub title\nParagraph" {
		t.Errorf("wrong Body: %s", post.Body)
	}
}
