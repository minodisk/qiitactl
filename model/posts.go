package model

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"text/template"
	"time"

	"github.com/briandowns/spinner"
)

const (
	postFileFormat = `<!--
{{.Meta.Format}}
-->
# {{.Title}}
{{.Body}}`
)

const (
	perPage = 100
)

var (
	tmpl = func() (t *template.Template) {
		t = template.New("postfile")
		template.Must(t.Parse(postFileFormat))
		return
	}()
	rInvalidFilename = regexp.MustCompile(`[^a-zA-Z0-9\-]+`)
	rHyphens         = regexp.MustCompile(`\-{2,}`)
)

type Posts []Post

func ShowPosts(client Client) (err error) {
	posts, err := FetchPosts(client)
	if err != nil {
		return
	}
	for _, post := range posts {
		fmt.Println(post.Id, post.CreatedAt.FormatDate(), post.Title)
	}
	return
}

func spin(ch chan bool) {
	s := spinner.New(spinner.CharSets[9], time.Millisecond*33)
	s.Start()
	for finished := range ch {
		log.Println(finished)
		if finished {
			s.Stop()
		}
	}
}

func FetchPosts(client Client) (posts Posts, err error) {
	return FetchPostsInTeam(client, Team{})
}

func FetchPostsInTeam(client Client, team Team) (posts Posts, err error) {
	v := url.Values{}
	v.Set("per_page", strconv.Itoa(perPage))
	s := spinner.New(spinner.CharSets[9], time.Millisecond*66)
	defer s.Stop()
	for page := 1; ; page++ {
		s.Stop()
		s.Prefix = fmt.Sprintf("Fetching posts from %d to %d: ", perPage*(page-1)+1, perPage*page)
		s.Start()
		v.Set("page", strconv.Itoa(page))
		body, err := client.Get(team.Name, "/authenticated_user/items", &v)
		if err != nil {
			return nil, err
		}
		var p Posts
		err = json.Unmarshal(body, &p)
		if err != nil {
			return nil, err
		}
		if len(p) == 0 {
			break
		}
		posts = append(posts, p...)
	}
	return
}

func (posts Posts) SaveToLocal(dirname string) (err error) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}

	dir := filepath.Join(wd, dirname)
	fmt.Printf("Make directory: %s\n", dir)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return
	}

	for _, post := range posts {
		path := filepath.Join(dir, post.generateFilename())
		fmt.Printf("Make file: %s\n", path)

		f, err := os.Create(path)
		defer f.Close()
		if err != nil {
			return err
		}
		err = tmpl.Execute(f, post)
		if err != nil {
			return err
		}
	}

	return
}
