package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"time"

	"github.com/minodisk/qiitactl/api"

	"gopkg.in/yaml.v2"
)

var (
	rPost = regexp.MustCompile(`^(?ms:\n*<!--(.*)-->\n*# +(.*?)\n+(.*))$`)
)

type Post struct {
	Meta
	User         User   `json:"user"`
	Title        string `json:"title"`         // 投稿のタイトル
	Body         string `json:"body"`          // Markdown形式の本文
	RenderedBody string `json:"rendered_body"` // HTML形式の本文
	Team         *Team  // チーム
	File         File   // ファイル
}

func (post *Post) UnmarshalJSON(data []byte) (err error) {
	type P Post
	var p P
	err = json.Unmarshal(data, &p)
	if err != nil {
		return
	}

	post.InitFile("")
	(*post) = Post(p)
	return
}

func NewPost() (post Post) {
	post.CreatedAt = Time{Time: time.Now()}
	post.UpdatedAt = post.CreatedAt
	post.InitFile("")
	return
}

func NewPostWithFile(path string) (post Post, err error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	post, err = NewPostWithBytes(b, path)
	return
}

func NewPostWithBytes(b []byte, path string) (post Post, err error) {
	matched := rPost.FindSubmatch(b)
	if len(matched) != 4 {
		err = fmt.Errorf("wrong format")
		return
	}

	err = yaml.Unmarshal((bytes.TrimSpace(matched[1])), &post.Meta)
	if err != nil {
		return
	}
	post.Title = string(bytes.TrimSpace(matched[2]))
	post.Body = string(bytes.TrimSpace(matched[3]))
	post.InitFile(path)
	return
}

func (post Post) InitFile(path string) {
	if path == "" {
		post.File.FillPath(post.CreatedAt, post.Title, post.Team)
	} else {
		post.File.Path = path
	}
}

func (post Post) BelongsToTeam() (b bool) {
	return post.Team != nil
}

func (post Post) Save() (err error) {
	err = post.File.Save(post)
	return
}

func (post Post) Create() (err error) {
	return
}

func (post Post) Update(client api.Client) (err error) {
	subDomain := ""
	if post.Team != nil {
		subDomain = post.Team.ID
	}
	body, err := client.Patch(subDomain, fmt.Sprintf("/items/%s", post.Id), post)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &post)
	if err != nil {
		return
	}
	err = post.File.Save(post)
	return
}

type Meta struct {
	Id        string `json:"id" yaml:"id"`                 // 投稿の一意なID
	Url       string `json:"url" yaml:"url"`               // 投稿のURL
	CreatedAt Time   `json:"created_at" yaml:"created_at"` // データが作成された日時
	UpdatedAt Time   `json:"updated_at" yaml:"updated_at"` // データが最後に更新された日時
	Private   bool   `json:"private" yaml:"private"`       // 限定共有状態かどうかを表すフラグ (Qiita:Teamでは無効)
	Coediting bool   `json:"coediting" yaml:"coediting"`   // この投稿が共同更新状態かどうか (Qiita:Teamでのみ有効)
	Tags      Tags   `json:"tags" yaml:"tags"`             // 投稿に付いたタグ一覧
}

func (meta Meta) Format() (out string) {
	o, err := yaml.Marshal(meta)
	if err != nil {
		panic(err)
	}
	out = string(bytes.TrimSpace(o))
	return
}

type Tag struct {
	Name     string   `json:"name"`
	Versions []string `json:"versions"`
}

type Tags []Tag

func (tags Tags) MarshalYAML() (data interface{}, err error) {
	obj := make(map[string][]string)
	for _, tag := range tags {
		obj[tag.Name] = tag.Versions
	}
	data = interface{}(obj)
	return
}

func (tags *Tags) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var t map[string][]string
	err = unmarshal(&t)
	if err != nil {
		return
	}

	for name, versions := range t {
		tag := Tag{
			Name:     name,
			Versions: versions,
		}
		*tags = append(*tags, tag)
	}

	return
}
