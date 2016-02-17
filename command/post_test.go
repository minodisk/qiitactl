package command_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/cli"
	"github.com/minodisk/qiitactl/model"
	"github.com/minodisk/qiitactl/testutil"
)

func TestMain(m *testing.M) {
	code := m.Run()
	// Clean up trashes
	testutil.CleanUp()
	os.Exit(code)
}

func TestFetchPostWithID(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	mux := http.NewServeMux()
	handleItem(mux)
	serverMine := httptest.NewServer(mux)
	defer serverMine.Close()
	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(func(subDomain, path string) (url string) {
		switch subDomain {
		case "":
			url = fmt.Sprintf("%s%s%s", serverMine.URL, "/api/v2", path)
		default:
			log.Fatalf("wrong sub domain \"%s\"", subDomain)
		}
		return
	})

	testutil.ShouldExistFile(t, 0)

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

# Example Title
## Example body`
	if actual != expected {
		t.Errorf("wrong body:\n%s", testutil.Diff(expected, actual))
	}
}

func TestFetchPostWithWrongID(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	mux := http.NewServeMux()
	handleItem(mux)
	serverMine := httptest.NewServer(mux)
	defer serverMine.Close()
	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(func(subDomain, path string) (url string) {
		switch subDomain {
		case "":
			url = fmt.Sprintf("%s%s%s", serverMine.URL, "/api/v2", path)
		default:
			log.Fatalf("wrong sub domain \"%s\"", subDomain)
		}
		return
	})

	testutil.ShouldExistFile(t, 0)

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

func TestFetchPostWithFilename(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	mux := http.NewServeMux()
	handleItem(mux)
	serverMine := httptest.NewServer(mux)
	defer serverMine.Close()
	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(func(subDomain, path string) (url string) {
		switch subDomain {
		case "":
			url = fmt.Sprintf("%s%s%s", serverMine.URL, "/api/v2", path)
		default:
			log.Fatalf("wrong sub domain \"%s\"", subDomain)
		}
		return
	})

	testutil.ShouldExistFile(t, 0)

	err = os.MkdirAll("mine/2000/01/01", 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("mine/2000/01/01/Example Title.md", []byte(`<!--
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

# Example old title
## Example old body`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	testutil.ShouldExistFile(t, 1)

	errBuf := bytes.NewBuffer([]byte{})
	app := cli.GenerateApp(client, os.Stdout, errBuf)
	err = app.Run([]string{"qiitactl", "fetch", "post", "-f", "mine/2000/01/01/Example Title.md"})
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

# Example Title
## Example body`
	if actual != expected {
		t.Errorf("wrong body:\n%s", testutil.Diff(expected, actual))
	}
}

func TestFetchPostWithWrongFilename(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	mux := http.NewServeMux()
	handleItem(mux)
	serverMine := httptest.NewServer(mux)
	defer serverMine.Close()
	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(func(subDomain, path string) (url string) {
		switch subDomain {
		case "":
			url = fmt.Sprintf("%s%s%s", serverMine.URL, "/api/v2", path)
		default:
			log.Fatalf("wrong sub domain \"%s\"", subDomain)
		}
		return
	})

	testutil.ShouldExistFile(t, 0)

	buf := bytes.NewBuffer([]byte{})
	errBuf := bytes.NewBuffer([]byte{})
	app := cli.GenerateApp(client, buf, errBuf)
	err = app.Run([]string{"qiitactl", "fetch", "post", "-f", "mine/2000/01/01/Example Title.md"})
	if err != nil {
		t.Fatal(err)
	}
	e := errBuf.Bytes()
	actual := string(e)
	expected := "open mine/2000/01/01/Example Title.md: no such file or directory"
	if actual != expected {
		t.Fatalf("error should occur when fetches post with wrong filename: %s", actual)
	}

	testutil.ShouldExistFile(t, 0)
}

func TestShowPostWithID(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	mux := http.NewServeMux()
	handleItem(mux)
	serverMine := httptest.NewServer(mux)
	defer serverMine.Close()
	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(func(subDomain, path string) (url string) {
		switch subDomain {
		case "":
			url = fmt.Sprintf("%s%s%s", serverMine.URL, "/api/v2", path)
		default:
			log.Fatalf("wrong sub domain \"%s\"", subDomain)
		}
		return
	})

	testutil.ShouldExistFile(t, 0)

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

	if string(buf.Bytes()) != `4bd431809afb1bb99e4f 2000/01/01 Example Title
` {
		t.Errorf("written text is wrong: %s", buf.Bytes())
	}
}

func TestShowPostWithWrongID(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	mux := http.NewServeMux()
	handleItem(mux)
	serverMine := httptest.NewServer(mux)
	defer serverMine.Close()
	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(func(subDomain, path string) (url string) {
		switch subDomain {
		case "":
			url = fmt.Sprintf("%s%s%s", serverMine.URL, "/api/v2", path)
		default:
			log.Fatalf("wrong sub domain \"%s\"", subDomain)
		}
		return
	})

	testutil.ShouldExistFile(t, 0)

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

func TestShowPostWithFilename(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	mux := http.NewServeMux()
	handleItem(mux)
	serverMine := httptest.NewServer(mux)
	defer serverMine.Close()
	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(func(subDomain, path string) (url string) {
		switch subDomain {
		case "":
			url = fmt.Sprintf("%s%s%s", serverMine.URL, "/api/v2", path)
		default:
			log.Fatalf("wrong sub domain \"%s\"", subDomain)
		}
		return
	})

	testutil.ShouldExistFile(t, 0)

	err = os.MkdirAll("mine/2000/01/01", 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("mine/2000/01/01/Example Title.md", []byte(`<!--
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

# Example old title
## Example old body`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	testutil.ShouldExistFile(t, 1)

	buf := bytes.NewBuffer([]byte{})
	errBuf := bytes.NewBuffer([]byte{})
	app := cli.GenerateApp(client, buf, errBuf)
	err = app.Run([]string{"qiitactl", "show", "post", "-f", "mine/2000/01/01/Example Title.md"})
	if err != nil {
		t.Fatal(err)
	}
	e := errBuf.Bytes()
	if len(e) != 0 {
		t.Fatal(string(e))
	}

	testutil.ShouldExistFile(t, 1)

	if string(buf.Bytes()) != `4bd431809afb1bb99e4f 2000/01/01 Example Title
` {
		t.Errorf("written text is wrong: %s", buf.Bytes())
	}
}

func TestShowPostWithWithWrongFilename(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	mux := http.NewServeMux()
	handleItem(mux)
	serverMine := httptest.NewServer(mux)
	defer serverMine.Close()
	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(func(subDomain, path string) (url string) {
		switch subDomain {
		case "":
			url = fmt.Sprintf("%s%s%s", serverMine.URL, "/api/v2", path)
		default:
			log.Fatalf("wrong sub domain \"%s\"", subDomain)
		}
		return
	})

	testutil.ShouldExistFile(t, 0)

	buf := bytes.NewBuffer([]byte{})
	errBuf := bytes.NewBuffer([]byte{})
	app := cli.GenerateApp(client, buf, errBuf)
	err = app.Run([]string{"qiitactl", "show", "post", "-f", "mine/2000/01/01/Example Title.md"})
	if err != nil {
		t.Fatal(err)
	}
	e := errBuf.Bytes()
	actual := string(e)
	expected := "open mine/2000/01/01/Example Title.md: no such file or directory"
	if actual != expected {
		t.Fatalf("error should occur when shows post with wrong filename: %s", actual)
	}

	testutil.ShouldExistFile(t, 0)
}

func TestShowPosts(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	mux := http.NewServeMux()
	handleAuthenticatedUserItems(mux)
	handleTeams(mux)
	serverMine := httptest.NewServer(mux)
	defer serverMine.Close()
	mux = http.NewServeMux()
	handleAuthenticatedUserItemsWithTeam(mux)
	serverTeam := httptest.NewServer(mux)
	defer serverTeam.Close()
	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(func(subDomain, path string) (url string) {
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

	testutil.ShouldExistFile(t, 0)

	buf := bytes.NewBuffer([]byte{})
	errBuf := bytes.NewBuffer([]byte{})
	app := cli.GenerateApp(client, buf, errBuf)
	err = app.Run([]string{"qiitactl", "show", "posts"})
	if err != nil {
		t.Fatal(err)
	}
	e := errBuf.Bytes()
	if len(e) != 0 {
		t.Fatal(string(e))
	}

	testutil.ShouldExistFile(t, 0)

	if string(buf.Bytes()) != `Posts in Qiita:
4bd431809afb1bb99e4f 2000/01/01 Example Title
Posts in Qiita:Team (Increments Inc.):
4bd431809afb1bb99e4t 2015/09/25 Example Title in team
` {
		t.Errorf("written text is wrong: %s", buf.Bytes())
	}
}

func TestShowPostsErrors(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	func() {
		mux := http.NewServeMux()
		// handleAuthenticatedUserItems(mux)
		handleTeams(mux)
		serverMine := httptest.NewServer(mux)
		defer serverMine.Close()
		mux = http.NewServeMux()
		handleAuthenticatedUserItemsWithTeam(mux)
		serverTeam := httptest.NewServer(mux)
		defer serverTeam.Close()
		err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
		if err != nil {
			log.Fatal(err)
		}
		client := api.NewClient(func(subDomain, path string) (url string) {
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

		testutil.ShouldExistFile(t, 0)

		buf := bytes.NewBuffer([]byte{})
		errBuf := bytes.NewBuffer([]byte{})
		app := cli.GenerateApp(client, buf, errBuf)
		err = app.Run([]string{"qiitactl", "show", "posts"})
		if err != nil {
			t.Fatal(err)
		}
		e := errBuf.Bytes()
		if len(e) == 0 {
			t.Fatal("error should occur")
		}
	}()

	func() {
		mux := http.NewServeMux()
		handleAuthenticatedUserItems(mux)
		// handleTeams(mux)
		serverMine := httptest.NewServer(mux)
		defer serverMine.Close()
		mux = http.NewServeMux()
		handleAuthenticatedUserItemsWithTeam(mux)
		serverTeam := httptest.NewServer(mux)
		defer serverTeam.Close()
		err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
		if err != nil {
			log.Fatal(err)
		}
		client := api.NewClient(func(subDomain, path string) (url string) {
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

		testutil.ShouldExistFile(t, 0)

		buf := bytes.NewBuffer([]byte{})
		errBuf := bytes.NewBuffer([]byte{})
		app := cli.GenerateApp(client, buf, errBuf)
		err = app.Run([]string{"qiitactl", "show", "posts"})
		if err != nil {
			t.Fatal(err)
		}
		e := errBuf.Bytes()
		if len(e) == 0 {
			t.Fatal("error should occur")
		}
	}()

	func() {
		mux := http.NewServeMux()
		handleAuthenticatedUserItems(mux)
		handleTeams(mux)
		serverMine := httptest.NewServer(mux)
		defer serverMine.Close()
		mux = http.NewServeMux()
		// handleAuthenticatedUserItemsWithTeam(mux)
		serverTeam := httptest.NewServer(mux)
		defer serverTeam.Close()
		err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
		if err != nil {
			log.Fatal(err)
		}
		client := api.NewClient(func(subDomain, path string) (url string) {
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

		testutil.ShouldExistFile(t, 0)

		buf := bytes.NewBuffer([]byte{})
		errBuf := bytes.NewBuffer([]byte{})
		app := cli.GenerateApp(client, buf, errBuf)
		err = app.Run([]string{"qiitactl", "show", "posts"})
		if err != nil {
			t.Fatal(err)
		}
		e := errBuf.Bytes()
		if len(e) == 0 {
			t.Fatal("error should occur")
		}
	}()
}

func TestFetchPosts(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	mux := http.NewServeMux()
	handleAuthenticatedUserItems(mux)
	handleTeams(mux)
	serverMine := httptest.NewServer(mux)
	defer serverMine.Close()
	mux = http.NewServeMux()
	handleAuthenticatedUserItemsWithTeam(mux)
	serverTeam := httptest.NewServer(mux)
	defer serverTeam.Close()
	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(func(subDomain, path string) (url string) {
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

	testutil.ShouldExistFile(t, 0)

	errBuf := bytes.NewBuffer([]byte{})
	app := cli.GenerateApp(client, os.Stdout, errBuf)
	err = app.Run([]string{"qiitactl", "fetch", "posts"})
	if err != nil {
		t.Fatal(err)
	}
	e := errBuf.Bytes()
	if len(e) != 0 {
		t.Fatal(string(e))
	}

	testutil.ShouldExistFile(t, 2)

	func() {
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

# Example Title
## Example body`
		if actual != expected {
			t.Errorf("wrong body:\n%s", testutil.Diff(expected, actual))
		}
	}()

	func() {
		b, err := ioutil.ReadFile("increments/2015/09/25/Example Title in team.md")
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
team:
  active: true
  id: increments
  name: Increments Inc.
-->

# Example Title in team
## Example body in team`
		if actual != expected {
			t.Errorf("wrong body:\n%s", testutil.Diff(expected, actual))
		}
	}()
}

func TestFetchPostsErrors(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	func() {
		mux := http.NewServeMux()
		// handleAuthenticatedUserItems(mux)
		handleTeams(mux)
		serverMine := httptest.NewServer(mux)
		defer serverMine.Close()
		mux = http.NewServeMux()
		handleAuthenticatedUserItemsWithTeam(mux)
		serverTeam := httptest.NewServer(mux)
		defer serverTeam.Close()
		err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
		if err != nil {
			log.Fatal(err)
		}
		client := api.NewClient(func(subDomain, path string) (url string) {
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

		buf := bytes.NewBuffer([]byte{})
		errBuf := bytes.NewBuffer([]byte{})
		app := cli.GenerateApp(client, buf, errBuf)
		err = app.Run([]string{"qiitactl", "fetch", "posts"})
		if err != nil {
			t.Fatal(err)
		}
		e := errBuf.Bytes()
		if len(e) == 0 {
			t.Fatal("error should occur")
		}
	}()

	func() {
		mux := http.NewServeMux()
		handleAuthenticatedUserItems(mux)
		// handleTeams(mux)
		serverMine := httptest.NewServer(mux)
		defer serverMine.Close()
		mux = http.NewServeMux()
		handleAuthenticatedUserItemsWithTeam(mux)
		serverTeam := httptest.NewServer(mux)
		defer serverTeam.Close()
		err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
		if err != nil {
			log.Fatal(err)
		}
		client := api.NewClient(func(subDomain, path string) (url string) {
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

		buf := bytes.NewBuffer([]byte{})
		errBuf := bytes.NewBuffer([]byte{})
		app := cli.GenerateApp(client, buf, errBuf)
		err = app.Run([]string{"qiitactl", "fetch", "posts"})
		if err != nil {
			t.Fatal(err)
		}
		e := errBuf.Bytes()
		if len(e) == 0 {
			t.Fatal("error should occur")
		}
	}()

	func() {
		mux := http.NewServeMux()
		handleAuthenticatedUserItems(mux)
		handleTeams(mux)
		serverMine := httptest.NewServer(mux)
		defer serverMine.Close()
		mux = http.NewServeMux()
		// handleAuthenticatedUserItemsWithTeam(mux)
		serverTeam := httptest.NewServer(mux)
		defer serverTeam.Close()
		err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
		if err != nil {
			log.Fatal(err)
		}
		client := api.NewClient(func(subDomain, path string) (url string) {
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

		buf := bytes.NewBuffer([]byte{})
		errBuf := bytes.NewBuffer([]byte{})
		app := cli.GenerateApp(client, buf, errBuf)
		err = app.Run([]string{"qiitactl", "fetch", "posts"})
		if err != nil {
			t.Fatal(err)
		}
		e := errBuf.Bytes()
		if len(e) == 0 {
			t.Fatal("error should occur")
		}
	}()
}

func TestCreatePost(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	mux := http.NewServeMux()
	handleItems(mux)
	serverMine := httptest.NewServer(mux)
	defer serverMine.Close()
	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(func(subDomain, path string) (url string) {
		switch subDomain {
		case "":
			url = fmt.Sprintf("%s%s%s", serverMine.URL, "/api/v2", path)
		default:
			log.Fatalf("wrong sub domain \"%s\"", subDomain)
		}
		return
	})

	testutil.ShouldExistFile(t, 0)

	err = os.MkdirAll("mine/2000/01/01", 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("mine/2000/01/01/Example Title.md", []byte(`<!--
id: ""
url: ""
created_at: 2000-01-01T09:00:00+09:00
updated_at: 2000-01-01T09:00:00+09:00
private: false
coediting: false
tags:
- Ruby:
  - 0.0.1
- NewTag:
  - "1.0"
team: null
-->

# Example Edited Title
## Example Edited Body`), 0664)
	if err != nil {
		t.Fatal(err)
	}

	testutil.ShouldExistFile(t, 1)

	errBuf := bytes.NewBuffer([]byte{})
	app := cli.GenerateApp(client, os.Stdout, errBuf)
	err = app.Run([]string{"qiitactl", "create", "post", "-f", "mine/2000/01/01/Example Title.md"})
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
created_at: 2016-02-01T21:51:42+09:00
updated_at: 2016-02-01T21:51:42+09:00
private: false
coediting: false
tags:
- Ruby:
  - 0.0.1
- NewTag:
  - "1.0"
team: null
-->

# Example Edited Title
## Example Edited Body`
	if actual != expected {
		t.Errorf("wrong content:\n%s", testutil.Diff(expected, actual))
	}
}

func TestCreatePostErrorWithNoFile(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	mux := http.NewServeMux()
	handleItems(mux)
	serverMine := httptest.NewServer(mux)
	defer serverMine.Close()
	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(func(subDomain, path string) (url string) {
		switch subDomain {
		case "":
			url = fmt.Sprintf("%s%s%s", serverMine.URL, "/api/v2", path)
		default:
			log.Fatalf("wrong sub domain \"%s\"", subDomain)
		}
		return
	})

	testutil.ShouldExistFile(t, 0)

	errBuf := bytes.NewBuffer([]byte{})
	app := cli.GenerateApp(client, os.Stdout, errBuf)
	err = app.Run([]string{"qiitactl", "create", "post", "-f", "mine/2000/01/01/Example Title.md"})
	if err != nil {
		t.Fatal(err)
	}
	e := errBuf.Bytes()
	if len(e) == 0 {
		t.Fatalf("error should occur")
	}
}

func TestCreatePostErrorWithNoServer(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	mux := http.NewServeMux()
	// handleItems(mux)
	serverMine := httptest.NewServer(mux)
	defer serverMine.Close()
	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(func(subDomain, path string) (url string) {
		switch subDomain {
		case "":
			url = fmt.Sprintf("%s%s%s", serverMine.URL, "/api/v2", path)
		default:
			log.Fatalf("wrong sub domain \"%s\"", subDomain)
		}
		return
	})

	testutil.ShouldExistFile(t, 0)

	err = os.MkdirAll("mine/2000/01/01", 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("mine/2000/01/01/Example Title.md", []byte(`<!--
id: ""
url: ""
created_at: 2000-01-01T09:00:00+09:00
updated_at: 2000-01-01T09:00:00+09:00
private: false
coediting: false
tags:
- Ruby:
  - 0.0.1
- NewTag:
  - "1.0"
team: null
-->

# Example Edited Title
## Example Edited Body`), 0664)
	if err != nil {
		t.Fatal(err)
	}

	testutil.ShouldExistFile(t, 1)

	errBuf := bytes.NewBuffer([]byte{})
	app := cli.GenerateApp(client, os.Stdout, errBuf)
	err = app.Run([]string{"qiitactl", "create", "post", "-f", "mine/2000/01/01/Example Title.md"})
	if err != nil {
		t.Fatal(err)
	}
	e := errBuf.Bytes()
	if len(e) == 0 {
		t.Fatal("error should occur")
	}

	testutil.ShouldExistFile(t, 1)
}

func TestUpdatePost(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	mux := http.NewServeMux()
	handleItem(mux)
	serverMine := httptest.NewServer(mux)
	defer serverMine.Close()
	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(func(subDomain, path string) (url string) {
		switch subDomain {
		case "":
			url = fmt.Sprintf("%s%s%s", serverMine.URL, "/api/v2", path)
		default:
			log.Fatalf("wrong sub domain \"%s\"", subDomain)
		}
		return
	})

	testutil.ShouldExistFile(t, 0)

	err = os.MkdirAll("mine/2000/01/01", 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("mine/2000/01/01/Example Title.md", []byte(`<!--
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
team: null
-->

# Example Edited Title
## Example Edited Body`), 0664)
	if err != nil {
		t.Fatal(err)
	}

	testutil.ShouldExistFile(t, 1)

	errBuf := bytes.NewBuffer([]byte{})
	app := cli.GenerateApp(client, os.Stdout, errBuf)
	err = app.Run([]string{"qiitactl", "update", "post", "-f", "mine/2000/01/01/Example Title.md"})
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
updated_at: 2016-02-01T21:51:42+09:00
private: false
coediting: false
tags:
- Ruby:
  - 0.0.1
- NewTag:
  - "1.0"
team: null
-->

# Example Edited Title
## Example Edited Body`
	if actual != expected {
		t.Errorf("wrong content:\n%s", testutil.Diff(expected, actual))
	}
}

func TestUpdatePostWithWrongFilename(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	mux := http.NewServeMux()
	handleItem(mux)
	serverMine := httptest.NewServer(mux)
	defer serverMine.Close()
	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(func(subDomain, path string) (url string) {
		switch subDomain {
		case "":
			url = fmt.Sprintf("%s%s%s", serverMine.URL, "/api/v2", path)
		default:
			log.Fatalf("wrong sub domain \"%s\"", subDomain)
		}
		return
	})

	testutil.ShouldExistFile(t, 0)

	buf := bytes.NewBuffer([]byte{})
	errBuf := bytes.NewBuffer([]byte{})
	app := cli.GenerateApp(client, buf, errBuf)
	err = app.Run([]string{"qiitactl", "update", "post", "-f", "mine/2000/01/01/Example Title.md"})
	if err != nil {
		t.Fatal(err)
	}
	e := errBuf.Bytes()
	actual := string(e)
	expected := "open mine/2000/01/01/Example Title.md: no such file or directory"
	if actual != expected {
		t.Fatalf("error should occur when updates post with wrong filename: %s", actual)
	}

	testutil.ShouldExistFile(t, 0)
}

func TestUpdatePostWithNoServer(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	mux := http.NewServeMux()
	// handleItem(mux)
	serverMine := httptest.NewServer(mux)
	defer serverMine.Close()
	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(func(subDomain, path string) (url string) {
		switch subDomain {
		case "":
			url = fmt.Sprintf("%s%s%s", serverMine.URL, "/api/v2", path)
		default:
			log.Fatalf("wrong sub domain \"%s\"", subDomain)
		}
		return
	})

	testutil.ShouldExistFile(t, 0)

	err = os.MkdirAll("mine/2000/01/01", 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("mine/2000/01/01/Example Title.md", []byte(`<!--
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
team: null
-->

# Example Edited Title
## Example Edited Body`), 0664)
	if err != nil {
		t.Fatal(err)
	}

	testutil.ShouldExistFile(t, 1)

	errBuf := bytes.NewBuffer([]byte{})
	app := cli.GenerateApp(client, os.Stdout, errBuf)
	err = app.Run([]string{"qiitactl", "update", "post", "-f", "mine/2000/01/01/Example Title.md"})
	if err != nil {
		t.Fatal(err)
	}
	e := errBuf.Bytes()
	if len(e) == 0 {
		t.Fatal("error should occur")
	}

	testutil.ShouldExistFile(t, 1)
}

func TestDeletePost(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	mux := http.NewServeMux()
	handleItem(mux)
	serverMine := httptest.NewServer(mux)
	defer serverMine.Close()
	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(func(subDomain, path string) (url string) {
		switch subDomain {
		case "":
			url = fmt.Sprintf("%s%s%s", serverMine.URL, "/api/v2", path)
		default:
			log.Fatalf("wrong sub domain \"%s\"", subDomain)
		}
		return
	})

	testutil.ShouldExistFile(t, 0)

	err = os.MkdirAll("mine/2000/01/01", 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("mine/2000/01/01/Example Title.md", []byte(`<!--
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
team: null
-->

# Example Edited Title
## Example Edited Body`), 0664)
	if err != nil {
		t.Fatal(err)
	}

	testutil.ShouldExistFile(t, 1)

	errBuf := bytes.NewBuffer([]byte{})
	app := cli.GenerateApp(client, os.Stdout, errBuf)
	err = app.Run([]string{"qiitactl", "delete", "post", "-f", "mine/2000/01/01/Example Title.md"})
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
- NewTag:
  - "1.0"
team: null
-->

# Example Edited Title
## Example Edited Body`
	if actual != expected {
		t.Errorf("wrong content:\n%s", testutil.Diff(expected, actual))
	}
}

func TestFetchDeleteWithWrongFilename(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	mux := http.NewServeMux()
	handleItem(mux)
	serverMine := httptest.NewServer(mux)
	defer serverMine.Close()
	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(func(subDomain, path string) (url string) {
		switch subDomain {
		case "":
			url = fmt.Sprintf("%s%s%s", serverMine.URL, "/api/v2", path)
		default:
			log.Fatalf("wrong sub domain \"%s\"", subDomain)
		}
		return
	})

	testutil.ShouldExistFile(t, 0)

	buf := bytes.NewBuffer([]byte{})
	errBuf := bytes.NewBuffer([]byte{})
	app := cli.GenerateApp(client, buf, errBuf)
	err = app.Run([]string{"qiitactl", "delete", "post", "-f", "mine/2000/01/01/Example Title.md"})
	if err != nil {
		t.Fatal(err)
	}
	e := errBuf.Bytes()
	actual := string(e)
	expected := "open mine/2000/01/01/Example Title.md: no such file or directory"
	if actual != expected {
		t.Fatalf("error should occur when deletes post with wrong filename: %s", actual)
	}

	testutil.ShouldExistFile(t, 0)
}

func TestDeletePostErrorWithNoServer(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	mux := http.NewServeMux()
	// handleItem(mux)
	serverMine := httptest.NewServer(mux)
	defer serverMine.Close()
	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(func(subDomain, path string) (url string) {
		switch subDomain {
		case "":
			url = fmt.Sprintf("%s%s%s", serverMine.URL, "/api/v2", path)
		default:
			log.Fatalf("wrong sub domain \"%s\"", subDomain)
		}
		return
	})

	testutil.ShouldExistFile(t, 0)

	err = os.MkdirAll("mine/2000/01/01", 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("mine/2000/01/01/Example Title.md", []byte(`<!--
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
team: null
-->

# Example Edited Title
## Example Edited Body`), 0664)
	if err != nil {
		t.Fatal(err)
	}

	testutil.ShouldExistFile(t, 1)

	errBuf := bytes.NewBuffer([]byte{})
	app := cli.GenerateApp(client, os.Stdout, errBuf)
	err = app.Run([]string{"qiitactl", "delete", "post", "-f", "mine/2000/01/01/Example Title.md"})
	if err != nil {
		t.Fatal(err)
	}
	e := errBuf.Bytes()
	if len(e) == 0 {
		t.Fatal("error should occur")
	}

	testutil.ShouldExistFile(t, 1)
}

func handleAuthenticatedUserItems(mux *http.ServeMux) {
	mux.HandleFunc("/api/v2/authenticated_user/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
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
					"title": "Example Title",
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
				testutil.ResponseError(w, 500, fmt.Errorf("shouldn't access over total count"))
				return
			}
			w.Header().Set("Total-Count", fmt.Sprint(1))
			w.Write([]byte(body))
		default:
			w.WriteHeader(405)
		}
	})
}

func handleTeams(mux *http.ServeMux) {
	mux.HandleFunc("/api/v2/teams", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Write([]byte(`[
			{
				"active": true,
				"id": "increments",
				"name": "Increments Inc."
			}
		]`))
		default:
			w.WriteHeader(405)
		}
	})
}

func handleItems(mux *http.ServeMux) {
	mux.HandleFunc("/api/v2/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			defer r.Body.Close()

			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				testutil.ResponseError(w, 500, err)
				return
			}
			if len(b) == 0 {
				testutil.ResponseAPIError(w, 500, api.ResponseError{
					Type:    "fatal",
					Message: "empty body",
				})
				return
			}

			type Options struct {
				Tweet *bool `json:"tweet"`
				Gist  *bool `json:"gist"`
			}
			var options Options
			err = json.Unmarshal(b, &options)
			if err != nil {
				testutil.ResponseError(w, 500, err)
				return
			}
			if options.Tweet == nil || options.Gist == nil {
				testutil.ResponseError(w, 500, errors.New("tweet or gist is required"))
				return
			}

			var post model.Post
			err = json.Unmarshal(b, &post)
			if err != nil {
				testutil.ResponseError(w, 500, err)
				return
			}
			post.ID = "4bd431809afb1bb99e4f"
			post.URL = "https://qiita.com/yaotti/items/4bd431809afb1bb99e4f"
			post.CreatedAt = model.Time{Time: time.Date(2016, 2, 1, 12, 51, 42, 0, time.UTC)}
			post.UpdatedAt = post.CreatedAt
			b, err = json.Marshal(post)
			if err != nil {
				testutil.ResponseError(w, 500, err)
				return
			}
			_, err = w.Write(b)
			if err != nil {
				testutil.ResponseError(w, 500, err)
				return
			}

		default:
			w.WriteHeader(405)
		}
	})
}

func handleItem(mux *http.ServeMux) {
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
					"title": "Example Title",
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
				testutil.ResponseError(w, 500, err)
				return
			}
			if string(b) == "" {
				testutil.ResponseAPIError(w, 500, api.ResponseError{
					Type:    "fatal",
					Message: "empty body",
				})
				return
			}

			var post model.Post
			err = json.Unmarshal(b, &post)
			if err != nil {
				testutil.ResponseError(w, 500, err)
				return
			}
			post.UpdatedAt = model.Time{Time: time.Date(2016, 2, 1, 12, 51, 42, 0, time.UTC)}
			b, err = json.Marshal(post)
			if err != nil {
				testutil.ResponseError(w, 500, err)
				return
			}
			_, err = w.Write(b)
			if err != nil {
				testutil.ResponseError(w, 500, err)
				return
			}

		case "DELETE":
			defer r.Body.Close()

			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				testutil.ResponseError(w, 500, err)
				return
			}
			if string(b) == "" {
				testutil.ResponseAPIError(w, 500, api.ResponseError{
					Type:    "fatal",
					Message: "empty body",
				})
				return
			}

			var post model.Post
			err = json.Unmarshal(b, &post)
			if err != nil {
				testutil.ResponseError(w, 500, err)
				return
			}
			b, err = json.Marshal(post)
			if err != nil {
				testutil.ResponseError(w, 500, err)
				return
			}
			_, err = w.Write(b)
			if err != nil {
				testutil.ResponseError(w, 500, err)
				return
			}

		default:
			w.WriteHeader(405)
		}
	})
}

func handleAuthenticatedUserItemsWithTeam(mux *http.ServeMux) {
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
					"title": "Example Title in team",
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
			testutil.ResponseError(w, 500, fmt.Errorf("shouldn't access over total count"))
			return
		}
		w.Header().Set("Total-Count", fmt.Sprint(1))
		w.Write([]byte(body))
	})
}
