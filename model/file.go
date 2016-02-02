package model

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	DirMine = "mine"
)

var (
	rInvalidFilename = regexp.MustCompile(`[^a-zA-Z0-9\-]+`)
	rHyphens         = regexp.MustCompile(`\-{2,}`)
)

type File struct {
	Path string
}

func (file *File) FillPath(createdAt ITime, title string, team *Team) {
	var dirname string
	if team == nil {
		dirname = DirMine
	} else {
		dirname = team.ID
	}
	dirname = filepath.Join(dirname, createdAt.Format("2006/01"))

	filename := fmt.Sprintf("%s-%s", createdAt.Format("02"), title)
	filename = rInvalidFilename.ReplaceAllString(filename, "-")
	filename = strings.ToLower(filename)
	filename = rHyphens.ReplaceAllString(filename, "-")
	filename = strings.TrimRight(filename, "-")
	filename = fmt.Sprintf("%s.md", filename)

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
	err = post.Encode(f)
	if err != nil {
		return
	}
	return
}
