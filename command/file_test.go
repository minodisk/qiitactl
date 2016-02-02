package command_test

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/minodisk/qiitactl/command"
	"github.com/minodisk/qiitactl/model"
)

func TestGenerateFile(t *testing.T) {
	defer func() {
		os.RemoveAll(model.DirMine)
	}()

	err := command.GenerateFile("", "Example Title")
	if err != nil {
		t.Fatal(err)
	}

	matched, err := filepath.Glob(fmt.Sprintf("%s/*/*/*.md", model.DirMine))
	if err != nil {
		t.Fatal(err)
	}
	if len(matched) != 1 {
		t.Fatalf("wrong number of files: %d", len(matched))
	}
	actual := matched[0]
	expected := fmt.Sprintf("%s/%s-example-title.md", model.DirMine, time.Now().Format("2006/01/02"))
	if actual != expected {
		t.Errorf("wrong path: expected %s, but actual %s", expected, actual)
	}
}

func TestGenerateUniqueFile(t *testing.T) {
	defer func() {
		os.RemoveAll(model.DirMine)
	}()

	var err error
	err = command.GenerateFile("", "Example Title")
	err = command.GenerateFile("", "Example Title")
	if err != nil {
		t.Fatal(err)
	}

	matched, err := filepath.Glob(fmt.Sprintf("%s/*/*/*.md", model.DirMine))
	if err != nil {
		t.Fatal(err)
	}

	sort.Strings(matched)

	func() {
		expected := 2
		actual := len(matched)
		if actual != expected {
			t.Fatalf("wrong number of files: expected %d, but actual %d", expected, actual)
		}
	}()

	func() {
		actual := matched[0]
		expected := fmt.Sprintf("%s/%s-example-title-.md", model.DirMine, time.Now().Format("2006/01/02"))
		if actual != expected {
			t.Errorf("wrong path: expected %s, but actual %s", expected, actual)
		}
	}()

	func() {
		actual := matched[1]
		expected := fmt.Sprintf("%s/%s-example-title.md", model.DirMine, time.Now().Format("2006/01/02"))
		if actual != expected {
			t.Errorf("wrong path: expected %s, but actual %s", expected, actual)
		}
	}()
}
