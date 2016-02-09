package testutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/minodisk/qiitactl/api"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func Diff(src1, src2 string) (s string) {
	dmp := diffmatchpatch.New()
	a, b, c := dmp.DiffLinesToChars(src1, src2)
	diffs := dmp.DiffMain(a, b, false)
	lines := dmp.DiffCharsToLines(diffs, c)

	var ls []string
	for _, line := range lines {
		switch line.Type {
		case diffmatchpatch.DiffDelete:
			ls = append(ls, prefix("- ", line.Text))
		case diffmatchpatch.DiffInsert:
			ls = append(ls, prefix("+ ", line.Text))
		}
	}
	s = strings.Join(ls, "\n")
	return
}

func prefix(pre, text string) (s string) {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = pre + line
	}
	s = strings.Join(lines, "\n")
	return
}

func CleanUp() {
	os.Unsetenv("QIITA_ACCESS_TOKEN")
	os.RemoveAll("mine")
	os.RemoveAll("increments")
	os.RemoveAll("foo")
}

func ResponseError(w http.ResponseWriter, statusCode int, err error) {

	ResponseAPIError(w, statusCode, api.ResponseError{
		Type:    "error",
		Message: err.Error(),
	})
}

func ResponseAPIError(w http.ResponseWriter, statusCode int, err api.ResponseError) {
	w.WriteHeader(statusCode)
	b, e := json.Marshal(err)
	if e != nil {
		fmt.Fprintf(w, "\"%s\"", e.Error())
		return
	}
	w.Write(b)
}

func ShouldExistFile(t *testing.T, num int) (paths []string) {
	paths, err := findMarkdownFiles()
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != num {
		t.Fatalf("file should exist %d file, but actual %s", num, paths)
	}
	return
}

func findMarkdownFiles() (pathes []string, err error) {
	err = filepath.Walk(".", func(p string, i os.FileInfo, e error) (err error) {
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
		pathes = append(pathes, p)
		return
	})
	return
}
