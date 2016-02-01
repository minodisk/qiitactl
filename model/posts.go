package model

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/minodisk/qiitactl/api"
)

const (
	perPage = 100
)

type Posts []Post

func FetchPosts(client api.Client, team *Team) (posts Posts, err error) {
	v := url.Values{}
	v.Set("per_page", strconv.Itoa(perPage))
	for page := 1; ; page++ {
		v.Set("page", strconv.Itoa(page))

		subDomain := ""
		if team != nil {
			subDomain = team.ID
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
	for i, post := range posts {
		post.Team = team
		post.InitFile("")
		posts[i] = post
	}
	return
}

func (posts Posts) Save() (err error) {
	for _, post := range posts {
		err = post.Save()
		if err != nil {
			return
		}
	}
	return
}
