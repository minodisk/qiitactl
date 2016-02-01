package model

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

const (
	DirMine        = "mine"
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

// type File struct {
// 	Post Post
// }
//
// func NewFile(post Post) (file File) {
// 	file.Post = post
// 	return
// }

type File struct {
	Path string
}

func (file *File) FillPath(createdAt Time, title string, team *Team) {
	filename := rInvalidFilename.ReplaceAllString(title, "-")
	filename = strings.ToLower(filename)
	filename = fmt.Sprintf("%s-%s", createdAt.ToPath(), filename)
	filename = rHyphens.ReplaceAllString(filename, "-")
	filename = strings.TrimRight(filename, "-")
	filename = fmt.Sprintf("%s.md", filename)
	var dirname string
	if team == nil {
		dirname = DirMine
	} else {
		dirname = team.ID
	}
	file.Path = filepath.Join(dirname, filename)
}

func (file *File) Save(post Post) (err error) {
	dir := filepath.Dir(file.Path)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return
	}

	fmt.Printf("Make file: %s\n", file.Path)
	f, err := os.Create(file.Path)
	defer f.Close()
	if err != nil {
		return
	}
	err = tmpl.Execute(f, post)
	if err != nil {
		return
	}
	return
}
