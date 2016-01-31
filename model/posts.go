package model

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"github.com/minodisk/qiitactl/api"
)

const (
	perPage = 100
)

type Posts []Post

func ShowPosts(client api.Client) (err error) {
	posts, err := FetchPosts(client)
	if err != nil {
		return
	}
	for _, post := range posts {
		fmt.Println(post.Id, post.CreatedAt.FormatDate(), post.Title)
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

func FetchPosts(client api.Client) (posts Posts, err error) {
	return FetchPostsInTeam(client, nil)
}

func FetchPostsInTeam(client api.Client, team *Team) (posts Posts, err error) {
	v := url.Values{}
	v.Set("per_page", strconv.Itoa(perPage))
	s := spinner.New(spinner.CharSets[9], time.Millisecond*66)
	defer s.Stop()
	for page := 1; ; page++ {
		s.Stop()
		s.Prefix = fmt.Sprintf("Fetching posts from %d to %d: ", perPage*(page-1)+1, perPage*page)
		s.Start()
		v.Set("page", strconv.Itoa(page))

		subDomain := ""
		if team != nil {
			subDomain = team.Name
		}
		body, err := client.Get(subDomain, "/authenticated_user/items", &v)
		if err != nil {
			return nil, err
		}

		var ps Posts
		err = json.Unmarshal(body, &ps)
		if err != nil {
			return nil, err
		}
		if len(ps) == 0 {
			break
		}
		posts = append(posts, ps...)
	}
	for _, post := range posts {
		post.Team = team
	}
	return
}

func (posts Posts) SaveToLocal() (err error) {
	for _, post := range posts {
		f := NewFile(post)
		err = f.Save()
		if err != nil {
			return
		}
	}
	return
}
