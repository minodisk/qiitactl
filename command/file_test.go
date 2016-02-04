package command_test

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/minodisk/qiitactl/cli"
	"github.com/minodisk/qiitactl/model"
	"github.com/minodisk/qiitactl/testutil"
)

func TestGenerateFile(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	app := cli.GenerateApp(client, os.Stdout, os.Stderr)
	err := app.Run([]string{"qiitactl", "generate", "file", "-t", "Example Title"})
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
	testutil.CleanUp()
	defer testutil.CleanUp()

	app := cli.GenerateApp(client, os.Stdout, os.Stderr)
	err := app.Run([]string{"qiitactl", "generate", "file", "-t", "Example Title"})
	if err != nil {
		t.Fatal(err)
	}
	err = app.Run([]string{"qiitactl", "generate", "file", "-t", "Example Title"})
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
