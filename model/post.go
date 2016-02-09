package model

import (
	"bytes"
	"encoding/json"
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

	// DirMine is the directory of saving posts in Qiita. (Not for posts in Qiita:Team)
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

// Post is a post in Qiita.
type Post struct {
	Meta
	User         User   `json:"user"`
	Title        string `json:"title"`         // 投稿のタイトル
	Body         string `json:"body"`          // Markdown形式の本文
	RenderedBody string `json:"rendered_body"` // HTML形式の本文
	Path         string `json:"-"`
}

// CreationOptions is options for creating a post.
type CreationOptions struct {
	Tweet bool `json:"tweet"`
	Gist  bool `json:"gist"`
}

// CreationPost is sent to Qiita server when creating post.
type CreationPost struct {
	Post
	CreationOptions
}

// Validate validates fields in Post.
func (post Post) Validate() (err InvalidError) {
	err = make(InvalidError)
	if post.Title == "" {
		err["title"] = InvalidStatus{
			Name:     "title",
			Required: true,
		}
	}
	if post.Body == "" {
		err["body"] = InvalidStatus{
			Name:     "body",
			Required: true,
		}
	}
	if post.Team == nil {
		if post.Tags == nil {
			err["tags"] = InvalidStatus{
				Name:     "tags",
				Required: true,
			}
		}
	}
	if err.none() {
		err = nil
	}
	return
}

// NewPost create a Post.
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

// NewPostWithFile loads local file and create a Post from the content of the file.
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

// Create creates a new post in Qiita.
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

// FetchPost fetches a post from Qiita.
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

// Update updates a post in Qiita.
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

// Delete deletes a post in Qiita.
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

// Save saves a post as a markdown file in local.
func (post *Post) Save(cachedPaths map[string]string) (err error) {
	if cachedPaths == nil {
		cachedPaths = pathsInLocal()
	}

	post.fillPath(cachedPaths)

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

func (post *Post) fillPath(paths map[string]string) {
	for id, path := range paths {
		if id == post.ID {
			post.Path = path
			return
		}
	}

	if post.Path != "" {
		return
	}

	post.Path = post.createPath()
	return
}

func pathsInLocal() (paths map[string]string) {
	paths = make(map[string]string)
	filepath.Walk(".", func(p string, i os.FileInfo, e error) (err error) {
		if e != nil {
			err = e
			return
		}
		if i.IsDir() {
			return
		}
		if filepath.Ext(p) != ".md" {
			return
		}

		post, err := NewPostWithFile(p)
		if err != nil {
			return nil
		}
		if post.ID == "" {
			return
		}
		paths[post.ID] = p
		return
	})
	return
}

func (post Post) createPath() (path string) {
	var dirname string
	if post.Team == nil {
		dirname = DirMine
	} else {
		dirname = post.Team.ID
	}
	dirname = filepath.Join(dirname, post.CreatedAt.Format("2006/01"))

	basename := fmt.Sprintf("%s-%s", post.CreatedAt.Format("02"), post.Title)
	basename = rInvalidFilename.ReplaceAllString(basename, "-")
	basename = strings.ToLower(basename)
	basename = rHyphens.ReplaceAllString(basename, "-")
	basename = strings.TrimRight(basename, "-")

	for {
		path = filepath.Join(dirname, fmt.Sprintf("%s.md", basename))
		_, err := os.Stat(path)
		// without error: file exists at the path
		// with error: file doesn't exist at the path
		if err != nil {
			// file doesn't exist at the path
			// so possible to create a file
			break
		}
		basename += "-"
	}
	return
}

// Encode encodes Post from bytes.
func (post Post) Encode(w io.Writer) (err error) {
	err = tmpl.Execute(w, post)
	return
}

// Decode decodes Post from bytes.
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

// EmptyIDError occurs when operate a post without ID.
type EmptyIDError struct{}

func (err EmptyIDError) Error() (msg string) {
	msg = "empty ID"
	return
}

// InvalidError occurs when some fields are wrong.
type InvalidError map[string]InvalidStatus

func (err InvalidError) Error() (msg string) {
	msgs := []string{}
	len := 0
	for _, status := range err {
		len++
		msgs = append(msgs, status.String())
	}
	var firstMessage string
	switch len {
	case 0:
		firstMessage = "Valid"
	case 1:
		firstMessage = "A field is invalid:"
	default:
		firstMessage = "Some fields are invalid:"
	}
	msgs = append([]string{firstMessage}, msgs...)
	msg = strings.Join(msgs, "\n")
	return
}

func (err InvalidError) none() (valid bool) {
	for range err {
		return false
	}
	return true
}

// InvalidStatus suggests what is wrong.
type InvalidStatus struct {
	Name     string
	Required bool
}

func (s InvalidStatus) String() (msg string) {
	msgs := []string{
		fmt.Sprintf("- %s", s.Name),
	}
	if s.Required {
		msgs = append(msgs, fmt.Sprintf("  - shouldn't be empty"))
	}
	msg = strings.Join(msgs, "\n")
	return
}
