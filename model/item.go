package model

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	rItem = regexp.MustCompile(`^(?ms:\n*<!--(.*)-->\n*# +(.*?)\n+(.*))$`)
)

type Item struct {
	Meta
	User         User   `json:"user"`
	Title        string `json:"title"`         // 投稿のタイトル
	Body         string `json:"body"`          // Markdown形式の本文
	RenderedBody string `json:"rendered_body"` // HTML形式の本文
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

func NewItemFromFile(filename string) (item Item, err error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	item, err = NewItem(body)
	return
}

func NewItem(b []byte) (item Item, err error) {
	matched := rItem.FindSubmatch(b)
	if len(matched) != 4 {
		err = fmt.Errorf("wrong format")
		return
	}
	item.Title = string(bytes.TrimSpace(matched[2]))
	item.Body = string(bytes.TrimSpace(matched[3]))
	err = yaml.Unmarshal((bytes.TrimSpace(matched[1])), &item.Meta)
	return
}

func (item Item) generateFilename() (f string) {
	b := fmt.Sprintf("%s-%s", item.CreatedAt.ToPath(), normalizeFilename(item.Title))
	f = fmt.Sprintf("%s.md", shortenHyphens(b))
	return
}

func normalizeFilename(filename string) (f string) {
	f = rInvalidFilename.ReplaceAllString(filename, "-")
	f = strings.ToLower(f)
	return
}

func shortenHyphens(filename string) (f string) {
	f = rHyphens.ReplaceAllString(filename, "-")
	f = strings.TrimRight(f, "-")
	return
}
