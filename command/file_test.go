package command_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/command"
	"github.com/minodisk/qiitactl/testutil"
)

func TestGenerateFileInMine(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(nil)

	testutil.ShouldExistFile(t, 0)

	app := command.New(client, os.Stdout, os.Stderr)
	app.Run([]string{"qiitactl", "generate", "file", "Example Title"})

	testutil.ShouldExistFile(t, 1)

	path := fmt.Sprintf("mine/%s/Example Title.md", time.Now().Format("2006/01/02"))
	_, err = os.Stat(path)
	if err != nil {
		t.Errorf("file should exist at %s", path)
	}
}

func TestGenerateFileInTeam(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(nil)

	testutil.ShouldExistFile(t, 0)

	app := command.New(client, os.Stdout, os.Stderr)
	app.Run([]string{"qiitactl", "generate", "file", "Example Title", "-t", "increments"})

	testutil.ShouldExistFile(t, 1)

	path := fmt.Sprintf("increments/%s/Example Title.md", time.Now().Format("2006/01/02"))
	_, err = os.Stat(path)
	if err != nil {
		t.Errorf("file should exist at %s", path)
	}
}

func TestGenerateUniqueFile(t *testing.T) {
	testutil.CleanUp()
	defer testutil.CleanUp()

	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client := api.NewClient(nil)

	testutil.ShouldExistFile(t, 0)

	app := command.New(client, os.Stdout, os.Stderr)
	app.Run([]string{"qiitactl", "generate", "file", "Example Title"})
	app.Run([]string{"qiitactl", "generate", "file", "Example Title"})

	testutil.ShouldExistFile(t, 2)

	path := fmt.Sprintf("mine/%s/Example Title.md", time.Now().Format("2006/01/02"))
	_, err = os.Stat(path)
	if err != nil {
		t.Errorf("file should exist at %s", path)
	}
	path = fmt.Sprintf("mine/%s/Example Title-.md", time.Now().Format("2006/01/02"))
	_, err = os.Stat(path)
	if err != nil {
		t.Errorf("file should exist at %s", path)
	}
}
