package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	BuildURL   func(string, string) string
	httpClient *http.Client
}

func BuildURL(subDomain, path string) (url string) {
	var host string
	if subDomain == "" {
		host = baseHost
	} else {
		host = fmt.Sprintf("%s.%s", subDomain, baseHost)
	}
	url = fmt.Sprintf("https://%s/api/v2%s", host, path)
	return
}

func NewClient(buildURL func(string, string) string) (c Client, err error) {
	if buildURL == nil {
		c.BuildURL = BuildURL
	} else {
		c.BuildURL = buildURL
	}

	c.Token = os.Getenv(envAccessToken)
	if c.Token == "" {
		err = fmt.Errorf("publish personal access token at https://qiita.com/settings/applications, then set environment variable as %s", envAccessToken)
		return
	}
	c.httpClient = &http.Client{}
	return
}

func (c Client) process(method string, subDomain string, path string, data interface{}) (respBody []byte, err error) {
	url := c.BuildURL(subDomain, path)
	// fmt.Println("->", method, url)

	var reqBody io.Reader
	if data != nil {
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(b)
	}
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	fmt.Println(resp.StatusCode, resp.Status)

	if resp.StatusCode/100 != 2 {
		e, err := NewError(respBody)
		if err != nil {
			return nil, err
		}
		err = e.Error()
		return nil, err
	}

	return
}

func (c Client) Post(subDomain string, path string, data interface{}) (body []byte, err error) {
	body, err = c.process("Post", subDomain, path, data)
	return
}

func (c Client) Get(subDomain string, path string, v *url.Values) (body []byte, err error) {
	if v != nil {
		path = fmt.Sprintf("%s?%s", path, v.Encode())
	}
	body, err = c.process("GET", subDomain, path, nil)
	return
}

func (c Client) Patch(subDomain string, path string, data interface{}) (body []byte, err error) {
	body, err = c.process("PATCH", subDomain, path, data)
	return
}

func (c Client) Delete(subDomain string, path string, data interface{}) (body []byte, err error) {
	body, err = c.process("Delete", subDomain, path, data)
	return
}
