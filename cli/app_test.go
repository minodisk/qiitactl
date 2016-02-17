package cli_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/cli"
	"github.com/minodisk/qiitactl/testutil"
)

var (
	server *httptest.Server
	client api.Client
)

func TestMain(m *testing.M) {
	mux := http.NewServeMux()
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

		default:
			w.WriteHeader(405)
		}
	})

	server = httptest.NewServer(mux)
	defer server.Close()

	client = api.NewClient(func(subDomain, path string) (url string) {
		switch subDomain {
		case "":
			url = fmt.Sprintf("%s%s%s", server.URL, "/api/v2", path)
		default:
			log.Fatalf("wrong sub domain \"%s\"", subDomain)
		}
		return
	})

	code := m.Run()
	os.Exit(code)
}

func TestNewAppNoCommand(t *testing.T) {
	outBuf := bytes.NewBuffer([]byte{})
	errBuf := bytes.NewBuffer([]byte{})
	app := cli.GenerateApp(client, outBuf, errBuf)
	err := app.Run([]string{"qiitactl"})
	if err != nil {
		t.Error(err)
	}
	if len(errBuf.Bytes()) != 0 {
		t.Fatalf("error shouldn't occur: %s", errBuf.Bytes())
	}
	out := string(outBuf.Bytes())
	if !strings.HasPrefix(out, "NAME:") {
		t.Errorf("wrong output: %s", out)
	}
}

func TestNewAppUndefinedCommand(t *testing.T) {
	outBuf := bytes.NewBuffer([]byte{})
	errBuf := bytes.NewBuffer([]byte{})
	app := cli.GenerateApp(client, outBuf, errBuf)
	err := app.Run([]string{"qiitactl", "nop"})
	if err != nil {
		t.Error(err)
	}
	if string(errBuf.Bytes()) != "qiitactl: 'nop' is not a qiitactl command. See 'qiitactl --help'.\n" {
		t.Fatalf("wrong error: %s", errBuf.Bytes())
	}
	if len(outBuf.Bytes()) != 0 {
		t.Errorf("shouldn't output: %s", outBuf.Bytes())
	}
}

func TestFetchPostWithID(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	testutil.ShouldExistFile(t, 0)

	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}

	errBuf := bytes.NewBuffer([]byte{})
	app := cli.GenerateApp(client, os.Stdout, errBuf)
	err = app.Run([]string{"qiitactl", "fetch", "post", "-i", "4bd431809afb1bb99e4f"})
	if err != nil {
		t.Fatal(err)
	}
	e := errBuf.Bytes()
	if len(e) != 0 {
		t.Fatal(string(e))
	}

	testutil.ShouldExistFile(t, 1)

	b, err := ioutil.ReadFile("mine/2000/01/01/Example Title.md")
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
team: null
-->

# Example title

## Example body`
	if actual != expected {
		t.Errorf("wrong body:\n%s", testutil.Diff(expected, actual))
	}
}

func TestFetchPostWithWrongID(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	testutil.ShouldExistFile(t, 0)

	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBuffer([]byte{})
	errBuf := bytes.NewBuffer([]byte{})
	app := cli.GenerateApp(client, buf, errBuf)
	err = app.Run([]string{"qiitactl", "fetch", "post", "-i", "XXXXXXXXXXXXXXXXXXXX"})
	if err != nil {
		t.Fatal(err)
	}
	e := errBuf.Bytes()
	actual := string(e)
	expected := "404 Not Found"
	if actual != expected {
		t.Fatalf("error should occur when fetches post with wrong ID: %s", actual)
	}

	testutil.ShouldExistFile(t, 0)
}

func TestShowPostWithID(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	testutil.ShouldExistFile(t, 0)

	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBuffer([]byte{})
	errBuf := bytes.NewBuffer([]byte{})
	app := cli.GenerateApp(client, buf, errBuf)
	err = app.Run([]string{"qiitactl", "show", "post", "-i", "4bd431809afb1bb99e4f"})
	if err != nil {
		t.Fatal(err)
	}
	e := errBuf.Bytes()
	if len(e) != 0 {
		t.Fatal(string(e))
	}

	testutil.ShouldExistFile(t, 0)

	if string(buf.Bytes()) != `4bd431809afb1bb99e4f 2000/01/01 Example title
` {
		t.Errorf("written text is wrong: %s", buf.Bytes())
	}
}

func TestShowPostWithWrongID(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	testutil.ShouldExistFile(t, 0)

	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBuffer([]byte{})
	errBuf := bytes.NewBuffer([]byte{})
	app := cli.GenerateApp(client, buf, errBuf)
	err = app.Run([]string{"qiitactl", "show", "post", "-i", "XXXXXXXXXXXXXXXXXXXX"})
	if err != nil {
		t.Fatal(err)
	}
	e := errBuf.Bytes()
	actual := string(e)
	expected := "404 Not Found"
	if actual != expected {
		t.Fatalf("error should occur when show post with wrong ID: %s", actual)
	}

	testutil.ShouldExistFile(t, 0)
}
