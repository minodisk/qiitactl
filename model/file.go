package model

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

const (
	postFileFormat = `<!--
{{.Meta.Format}}
-->
# {{.Title}}
{{.Body}}`
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

type File struct {
	Post Post
}

func NewFile(post Post) (file File) {
	file.Post = post
	return
}

func NewFileFromLocal(filename string) (file File, err error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	file.Post, err = NewPostWithBytes(b)
	return
}

func (file *File) Save() (err error) {
	var dir string
	if file.Post.BelongsToTeam() {
		dir = file.Post.Team.Name
	} else {
		dir = "mine"
	}

	fmt.Printf("Make directory: %s\n", dir)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return
	}

	path := filepath.Join(dir, generateFilename(file.Post))
	fmt.Printf("Make file: %s\n", path)
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return
	}
	err = tmpl.Execute(f, file.Post)
	if err != nil {
		return
	}
	return
}

func generateFilename(post Post) (f string) {
	b := fmt.Sprintf("%s-%s", post.CreatedAt.ToPath(), normalizeFilename(post.Title))
	f = fmt.Sprintf("%s.md", shortenHyphens(b))
	return
}

func normalizeFilename(filename string) (f string) {
	f = rInvalidFilename.ReplaceAllString(filename, "-")
	f = strings.ToLower(f)
	return
}

func shortenHyphens(filename string) (f string) {
	f = rHyphens.ReplaceAllString(filename, "-")
	f = strings.TrimRight(f, "-")
	return
}
