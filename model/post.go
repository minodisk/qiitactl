package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/minodisk/qiitactl/api"

	"gopkg.in/yaml.v2"
)

const (
	postTemplate = `<!--
{{.Meta.Encode}}
-->
# {{.Title}}
{{.Body}}`
	DirMine = "mine"
)

var (
	rPostDecoder = regexp.MustCompile(`^(?ms:\n*<!--(.*)-->\n*# +(.*?)\n+(.*))$`)
	tmpl         = func() (t *template.Template) {
		t = template.New("postfile")
		template.Must(t.Parse(postTemplate))
		return
	}()
	rInvalidFilename = regexp.MustCompile(`[^a-zA-Z0-9\-]+`)
	rHyphens         = regexp.MustCompile(`\-{2,}`)
)

type Post struct {
	Meta
	User         User   `json:"user"`
	Title        string `json:"title"`         // 投稿のタイトル
	Body         string `json:"body"`          // Markdown形式の本文
	RenderedBody string `json:"rendered_body"` // HTML形式の本文
	Path         string `json:"-"`
}

type CreationOptions struct {
	Tweet bool `json:"tweet"`
	Gist  bool `json:"gist"`
}

type CreationPost struct {
	Post
	CreationOptions
}

func NewPost(title string, createdAt *Time, team *Team) (post Post) {
	if createdAt == nil {
		createdAt = &Time{Time: time.Now()}
	}
	post.CreatedAt = *createdAt
	post.UpdatedAt = *createdAt
	post.Title = title
	post.Team = team
	return
}

func NewPostWithFile(path string) (post Post, err error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = post.Decode(b)
	if err != nil {
		return
	}
	post.Path = path
	return
}

func (post *Post) Create(client api.Client, opts CreationOptions) (err error) {
	subDomain := ""
	if post.Team != nil {
		subDomain = post.Team.ID
	}

	cPost := CreationPost{
		Post:            *post,
		CreationOptions: opts,
	}
	body, _, err := client.Post(subDomain, "/items", cPost)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, post)
	if err != nil {
		return
	}
	return
}

func FetchPost(client api.Client, team *Team, id string) (post Post, err error) {
	subDomain := ""
	if team != nil {
		subDomain = team.ID
	}
	body, _, err := client.Get(subDomain, fmt.Sprintf("/items/%s", id), nil)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &post)
	if err != nil {
		return
	}
	post.Team = team
	return
}

func (post *Post) Update(client api.Client) (err error) {
	if post.ID == "" {
		err = EmptyIDError{}
		return
	}

	subDomain := ""
	if post.Team != nil {
		subDomain = post.Team.ID
	}
	body, _, err := client.Patch(subDomain, fmt.Sprintf("/items/%s", post.ID), post)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, post)
	if err != nil {
		return
	}
	return
}

func (post *Post) Delete(client api.Client) (err error) {
	if post.ID == "" {
		err = EmptyIDError{}
		return
	}

	subDomain := ""
	if post.Team != nil {
		subDomain = post.Team.ID
	}
	body, _, err := client.Delete(subDomain, fmt.Sprintf("/items/%s", post.ID), post)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, post)
	if err != nil {
		return
	}
	return
}

func (post *Post) Save() (err error) {
	err = post.FillPath()
	if err != nil {
		return
	}

	dir := filepath.Dir(post.Path)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return
	}

	fmt.Printf("Make file: %s\n", post.Path)

	f, err := os.Create(post.Path)
	defer f.Close()
	if err != nil {
		return
	}
	err = post.Encode(f)
	if err != nil {
		return
	}
	return
}

func (post *Post) FillPath() (err error) {
	if post.Path != "" {
		return
	}

	if post.ID != "" {
		path, err := post.findPath()
		if err == nil {
			post.Path = path
			return err
		}
	}

	post.Path = post.createPath()
	return
}

func (post Post) findPath() (path string, err error) {
	found := errors.New("found")
	err = filepath.Walk(".", func(p string, info os.FileInfo, e error) (err error) {
		if e != nil {
			err = e
			return
		}
		if info.IsDir() {
			return
		}
		if filepath.Ext(p) != ".md" {
			return
		}

		postInLocal, err := NewPostWithFile(p)
		if err != nil {
			err = nil
			return
		}
		if postInLocal.ID == post.ID {
			path = p
			return found
		}
		return
	})
	if err != found {
		err = PathNotFoundError{}
		return
	}
	return path, nil
}

func (post Post) createPath() (path string) {
	var dirname string
	if post.Team == nil {
		dirname = DirMine
	} else {
		dirname = post.Team.ID
	}
	dirname = filepath.Join(dirname, post.CreatedAt.Format("2006/01"))

	filename := fmt.Sprintf("%s-%s", post.CreatedAt.Format("02"), post.Title)
	filename = rInvalidFilename.ReplaceAllString(filename, "-")
	filename = strings.ToLower(filename)
	filename = rHyphens.ReplaceAllString(filename, "-")
	filename = strings.TrimRight(filename, "-")

	for {
		path = filepath.Join(dirname, fmt.Sprintf("%s.md", filename))
		_, err := os.Stat(path)
		// no error means: a file exists at the path
		// error occurs means: no file exists at the path
		if err != nil {
			// no file at the path,
			// so possible to create file with the path
			break
		}
		filename += "-"
	}
	return
}

func (post Post) Encode(w io.Writer) (err error) {
	err = tmpl.Execute(w, post)
	return
}

func (post *Post) Decode(b []byte) (err error) {
	matched := rPostDecoder.FindSubmatch(b)
	if len(matched) != 4 {
		err = fmt.Errorf("wrong format")
		return
	}

	err = yaml.Unmarshal((bytes.TrimSpace(matched[1])), &post.Meta)
	if err != nil {
		return
	}
	post.Title = string(bytes.TrimSpace(matched[2]))
	post.Body = string(bytes.TrimSpace(matched[3]))
	return
}

type EmptyIDError struct{}

func (err EmptyIDError) Error() (msg string) {
	msg = "empty ID"
	return
}

type PathNotFoundError struct{}

func (err PathNotFoundError) Error() (msg string) {
	msg = "path not found"
	return
}
