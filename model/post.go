package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"text/template"
	"time"

	"github.com/minodisk/qiitactl/api"

	"gopkg.in/yaml.v2"
)

const (
	postTemplate = `<!--
{{.Meta}}
-->
# {{.Title}}
{{.Body}}`
)

var (
	rPostDecoder = regexp.MustCompile(`^(?ms:\n*<!--(.*)-->\n*# +(.*?)\n+(.*))$`)
	tmpl         = func() (t *template.Template) {
		t = template.New("postfile")
		template.Must(t.Parse(postTemplate))
		return
	}()
)

type Post struct {
	Meta
	User         User   `json:"user"`
	Title        string `json:"title"`         // 投稿のタイトル
	Body         string `json:"body"`          // Markdown形式の本文
	RenderedBody string `json:"rendered_body"` // HTML形式の本文
	Team         *Team  // チーム
	File         File   // ファイル
}

func (post *Post) UnmarshalJSON(data []byte) (err error) {
	type P Post
	var p P
	err = json.Unmarshal(data, &p)
	if err != nil {
		return
	}

	post.FillFilePath()
	(*post) = Post(p)
	return
}

func NewPost(title string, team *Team) (post Post) {
	post.CreatedAt = Time{Time: time.Now()}
	post.UpdatedAt = post.CreatedAt
	post.Title = title
	post.Team = team
	post.FillFilePath()
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
	post.File.Path = path
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

func (post *Post) FillFilePath() {
	post.File.FillPath(post.CreatedAt, post.Title, post.Team)
}

func (post Post) Save() (err error) {
	err = post.File.Save(post)
	return
}

func (post Post) Create() (err error) {
	return
}

func (post Post) Update(client api.Client) (err error) {
	subDomain := ""
	if post.Team != nil {
		subDomain = post.Team.ID
	}
	body, err := client.Patch(subDomain, fmt.Sprintf("/items/%s", post.Id), post)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &post)
	if err != nil {
		return
	}
	err = post.File.Save(post)
	return
}
