package testutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

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
	os.RemoveAll("mine")
	os.RemoveAll("increments")
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
