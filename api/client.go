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
	"strings"

	"github.com/fatih/color"
	"github.com/minodisk/qiitactl/info"
)

const (
	envAccessToken = "QIITA_ACCESS_TOKEN"
	baseHost       = "qiita.com"
)

type Client struct {
	Token      string
	BuildURL   func(string, string) string
	httpClient *http.Client
	debug      bool
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
		err = EmptyTokenError{}
		return
	}
	c.httpClient = &http.Client{}
	return
}

func (c *Client) SetDebug(debug bool) {
	c.debug = debug
}

func (c Client) Process(method string, subDomain string, path string, data interface{}) (respBody []byte, respHeader http.Header, err error) {
	url := c.BuildURL(subDomain, path)

	var reqBody io.Reader
	if data != nil {
		b, err := json.Marshal(data)
		if err != nil {
			return nil, nil, err
		}
		reqBody = bytes.NewBuffer(b)
	}
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return
	}
	req.Header.Add("User-Agent", fmt.Sprintf("%s/%s", info.Name, info.Version))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	if data != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	if c.debug {
		blueBold := color.New(color.FgBlue).SprintFunc()
		blue := color.New(color.FgBlue).SprintFunc()
		magenta := color.New(color.FgCyan).SprintFunc()
		white := color.New(color.FgWhite).SprintFunc()
		fmt.Printf("%s %s %s\n%s\n%s\n", blueBold(req.Method), blue(req.URL), blue(req.Proto), magenta(stringifyHeader(req.Header)), white(stringifyBody(data)))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}

	respHeader = resp.Header

	defer resp.Body.Close()
	respBody, err = ioutil.ReadAll(resp.Body)
	if c.debug {
		blue := color.New(color.FgBlue).SprintFunc()
		magenta := color.New(color.FgCyan).SprintFunc()
		white := color.New(color.FgWhite).SprintFunc()
		fmt.Printf("%s %s\n%s\n%s\n", blue(resp.Proto), blue(resp.StatusCode), magenta(stringifyHeader(respHeader)), white(string(respBody)))
	}
	if err != nil {
		return
	}

	if resp.StatusCode/100 == 2 {
		return
	}

	switch resp.StatusCode {
	case 401:
		err = WrongTokenError{}
		return
	default:
		var respError ResponseError
		err = json.Unmarshal(respBody, &respError)
		if err == nil {
			err = respError
			return
		}
		err = StatusError{
			Code:    resp.StatusCode,
			Message: resp.Status,
		}
		return
	}

	return
}

func stringifyHeader(header http.Header) string {
	var lines []string
	for key, val := range header {
		for _, v := range val {
			lines = append(lines, fmt.Sprintf("%s: %s", key, v))
		}
	}
	return strings.Join(lines, "\n")
}

func stringifyBody(data interface{}) string {
	if data == nil {
		return ""
	}
	str, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return ""
	}
	return string(str)
}

func (c Client) Post(subDomain string, path string, data interface{}) (body []byte, header http.Header, err error) {
	body, header, err = c.Process("POST", subDomain, path, data)
	return
}

func (c Client) Get(subDomain string, path string, v *url.Values) (body []byte, header http.Header, err error) {
	if v != nil {
		path = fmt.Sprintf("%s?%s", path, v.Encode())
	}
	body, header, err = c.Process("GET", subDomain, path, nil)
	return
}

func (c Client) Patch(subDomain string, path string, data interface{}) (body []byte, header http.Header, err error) {
	body, header, err = c.Process("PATCH", subDomain, path, data)
	return
}

func (c Client) Delete(subDomain string, path string, data interface{}) (body []byte, header http.Header, err error) {
	body, header, err = c.Process("DELETE", subDomain, path, data)
	return
}

type EmptyTokenError struct{}

func (err EmptyTokenError) Error() (msg string) {
	msg = fmt.Sprintf("empty token: publish personal access token at https://qiita.com/settings/applications, then set environment variable as %s", envAccessToken)
	return
}

type WrongTokenError struct{}

func (err WrongTokenError) Error() (msg string) {
	msg = fmt.Sprintf("wrong token: publish personal access token at https://qiita.com/settings/applications, then set environment variable as %s", envAccessToken)
	return
}

type ResponseError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (err ResponseError) Error() (msg string) {
	msg = err.Message
	return
}

type StatusError struct {
	Code    int
	Message string
}

func (err StatusError) Error() (msg string) {
	msg = err.Message
	return
}
