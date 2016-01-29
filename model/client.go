package model

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const (
	baseHost = "qiita.com"
)

type Client struct {
	Token      string
	httpClient *http.Client
}

func NewClient() (c Client) {
	c = Client{
		Token:      os.Getenv("QIITA_ACCESS_TOKEN"),
		httpClient: &http.Client{},
	}
	return
}

func (c Client) process(method string, subDomain string, p string) (body []byte, err error) {
	var host string
	if subDomain == "" {
		host = baseHost
	} else {
		host = fmt.Sprintf("%s.%s", subDomain, baseHost)
	}
	url := fmt.Sprintf("https://%s/api/v2%s", host, p)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}

func (c Client) Get(subDomain string, p string, v *url.Values) (body []byte, err error) {
	body, err = c.process("GET", subDomain, fmt.Sprintf("%s?%s", p, v.Encode()))
	return
}
