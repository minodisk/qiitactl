package model

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"text/template"
	"time"

	"github.com/briandowns/spinner"
	"github.com/minodisk/qiitactl/api"
)

const (
	itemFileFormat = `<!--
{{.Meta.Format}}
-->
# {{.Title}}
{{.Body}}`
)

const (
	perPage = 100
)

var (
	tmpl = func() (t *template.Template) {
		t = template.New("itemfile")
		template.Must(t.Parse(itemFileFormat))
		return
	}()
	rInvalidFilename = regexp.MustCompile(`[^a-zA-Z0-9\-]+`)
	rHyphens         = regexp.MustCompile(`\-{2,}`)
)

type Items []Item

func ShowItems(client api.Client) (err error) {
	items, err := FetchItems(client)
	if err != nil {
		return
	}
	for _, item := range items {
		fmt.Println(item.Id, item.CreatedAt.FormatDate(), item.Title)
	}
	return
}

func spin(ch chan bool) {
	s := spinner.New(spinner.CharSets[9], time.Millisecond*33)
	s.Start()
	for finished := range ch {
		log.Println(finished)
		if finished {
			s.Stop()
		}
	}
}

func FetchItems(client api.Client) (items Items, err error) {
	return FetchItemsInTeam(client, Team{})
}

func FetchItemsInTeam(client api.Client, team Team) (items Items, err error) {
	v := url.Values{}
	v.Set("per_page", strconv.Itoa(perPage))
	s := spinner.New(spinner.CharSets[9], time.Millisecond*66)
	defer s.Stop()
	for page := 1; ; page++ {
		s.Stop()
		s.Prefix = fmt.Sprintf("Fetching items from %d to %d: ", perPage*(page-1)+1, perPage*page)
		s.Start()
		v.Set("page", strconv.Itoa(page))
		body, err := client.Get(team.Name, "/authenticated_user/items", &v)
		if err != nil {
			return nil, err
		}
		var p Items
		err = json.Unmarshal(body, &p)
		if err != nil {
			return nil, err
		}
		if len(p) == 0 {
			break
		}
		items = append(items, p...)
	}
	return
}

func (items Items) SaveToLocal(dirname string) (err error) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}

	dir := filepath.Join(wd, dirname)
	fmt.Printf("Make directory: %s\n", dir)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return
	}

	for _, item := range items {
		path := filepath.Join(dir, item.generateFilename())
		fmt.Printf("Make file: %s\n", path)

		f, err := os.Create(path)
		defer f.Close()
		if err != nil {
			return err
		}
		err = tmpl.Execute(f, item)
		if err != nil {
			return err
		}
	}

	return
}
