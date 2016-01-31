package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const (
	envAccessToken = "QIITA_ACCESS_TOKEN"
	baseHost       = "qiita.com"
)

type Client struct {
	Token      string
	httpClient *http.Client
}

func NewClient() (c Client, err error) {
	c.Token = os.Getenv(envAccessToken)
	if c.Token == "" {
		err = fmt.Errorf("publish personal access token at https://qiita.com/settings/applications, then set environment variable as %s", envAccessToken)
		return
	}
	c.httpClient = &http.Client{}
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
	// fmt.Println(method, url)
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
	if err != nil {
		return
	}

	if resp.StatusCode/100 != 2 {
		e, err := NewError(body)
		if err != nil {
			return nil, err
		}
		err = e.Error()
		return nil, err
	}

	return
}

func (c Client) Get(subDomain string, p string, v *url.Values) (body []byte, err error) {
	if v != nil {
		p = fmt.Sprintf("%s?%s", p, v.Encode())
	}
	body, err = c.process("GET", subDomain, p)
	return
}
