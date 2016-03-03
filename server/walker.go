package server

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
)

type Element struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Rel      string     `json:"rel"`
	Abs      string     `json:"abs"`
	Children []*Element `json:"children"`
}

func NewElement(path, rootPath, name string) (el Element, err error) {
	el.Abs = path
	el.Rel, err = filepath.Rel(rootPath, path)
	if err != nil {
		return
	}
	hasher := md5.New()
	hasher.Write([]byte(el.Abs))
	el.ID = hex.EncodeToString(hasher.Sum(nil))
	if name == "" {
		names := strings.Split(path, "/")
		el.Name = names[len(names)-1]
	} else {
		el.Name = name
	}
	return
}

func getTree() (root Element, err error) {
	rootPath, err := os.Getwd()
	if err != nil {
		return
	}
	root, err = NewElement(rootPath, rootPath, "")
	if err != nil {
		return
	}

	var paths []string

	err = filepath.Walk(rootPath, func(path string, i os.FileInfo, e error) (err error) {
		if e != nil {
			err = e
			return
		}
		if i.IsDir() {
			return
		}
		if filepath.Ext(path) != ".md" {
			return
		}

		path, err = filepath.Rel(root.Abs, path)
		if err != nil {
			return
		}
		paths = append(paths, path)
		return
	})

	for _, path := range paths {
		names := strings.Split(path, "/")
		p := root.Abs
		parent := &root
		for _, name := range names {
			p = filepath.Join(p, name)
			found := false
			var child *Element
			for _, child = range parent.Children {
				if child.Name == name {
					found = true
					break
				}
			}
			if !found {
				el, err := NewElement(p, rootPath, name)
				if err != nil {
					return root, err
				}
				child = &el
				parent.Children = append(parent.Children, child)
			}
			parent = child
		}
	}

	return
}
