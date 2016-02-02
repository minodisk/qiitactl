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

func TestFetchPosts(t *testing.T) {
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

	server := httptest.NewServer(mux)
	defer server.Close()

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

	posts, err := model.FetchPosts(client, nil)
	if err != nil {
		log.Fatal(err)
	}

	if len(posts) != 1 {
		log.Fatalf("wrong length: expected %d, but actual %d", 1, len(posts))
	}

	post := posts[0]
	if post.RenderedBody != "<h2>Example body</h2>" {
		log.Fatalf("wrong RenderedBody")
	}
	if post.Body != "## Example body" {
		log.Fatalf("wrong Body")
	}
	if post.Coediting != false {
		log.Fatalf("wrong Coediting")
	}
	// if post.CreatedAt != "2000-01-01T00:00:00+00:00" {
	// 	log.Fatalf("wrong CreatedAt")
	// }
	if post.ID != "4bd431809afb1bb99e4f" {
		log.Fatalf("wrong ID")
	}
	if post.Private != false {
		log.Fatalf("wrong Private")
	}
	// "tags": [
	// 	{
	// 		"name": "Ruby",
	// 		"versions": [
	// 			"0.0.1"
	// 		]
	// 	}
	// ],
	// "title": "Example title",
	// "updated_at": "2000-01-01T00:00:00+00:00",
	// "url": "https://qiita.com/yaotti/items/4bd431809afb1bb99e4f",
	// "user": {
	// 	"description": "Hello, world.",
	// 	"facebook_id": "yaotti",
	// 	"followees_count": 100,
	// 	"followers_count": 200,
	// 	"github_login_name": "yaotti",
	// 	"id": "yaotti",
	// 	"items_count": 300,
	// 	"linkedin_id": "yaotti",
	// 	"location": "Tokyo, Japan",
	// 	"name": "Hiroshige Umino",
	// 	"organization": "Increments Inc",
	// 	"permanent_id": 1,
	// 	"profile_image_url": "https://si0.twimg.com/profile_images/2309761038/1ijg13pfs0dg84sk2y0h_normal.jpeg",
	// 	"twitter_screen_name": "yaotti",
	// 	"website_url": "http://yaotti.hatenablog.com"
	// }
}
