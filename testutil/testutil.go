package testutil

import (
	"os"
	"strings"

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
