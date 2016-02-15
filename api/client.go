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

// Client is HTTP client accessing to the Qiita API v2.
type Client struct {
	BuildURL   func(string, string) string
	httpClient *http.Client
	debugMode  bool
}

// BuildURL builds URL of Qiita API v2.
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

// NewClient makes a Client.
// Client will access to URL made by buildURL.
func NewClient(buildURL func(string, string) string) (c Client) {
	if buildURL == nil {
		c.BuildURL = BuildURL
	} else {
		c.BuildURL = buildURL
	}

	c.httpClient = &http.Client{}
	return
}

// DebugMode sets debugMode to Client.
// When debugMode is true, qiitactl outputs the logs (e.g., HTTP request and response).
func (c *Client) DebugMode(debugMode bool) {
	c.debugMode = debugMode
}

func (c Client) process(method string, subDomain string, path string, data interface{}) (respBody []byte, respHeader http.Header, err error) {
	token := os.Getenv(envAccessToken)
	if token == "" {
		err = EmptyTokenError{}
		return
	}

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
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	if data != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	if c.debugMode {
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
	if c.debugMode {
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

// Options send OPTIONS request with data body
// to the URL built with subDomain and path.
func (c Client) Options(subDomain string, path string, data interface{}) (body []byte, header http.Header, err error) {
	body, header, err = c.process("OPTIONS", subDomain, path, data)
	return
}

// Post send POST request with data body
// to the URL built with subDomain and path.
func (c Client) Post(subDomain string, path string, data interface{}) (body []byte, header http.Header, err error) {
	body, header, err = c.process("POST", subDomain, path, data)
	return
}

// Get send GET request
// to the URL built with subDomain and path.
func (c Client) Get(subDomain string, path string, v *url.Values) (body []byte, header http.Header, err error) {
	if v != nil {
		path = fmt.Sprintf("%s?%s", path, v.Encode())
	}
	body, header, err = c.process("GET", subDomain, path, nil)
	return
}

// Patch send PATCH request with data body
// to the URL built with subDomain and path.
func (c Client) Patch(subDomain string, path string, data interface{}) (body []byte, header http.Header, err error) {
	body, header, err = c.process("PATCH", subDomain, path, data)
	return
}

// Delete send DELETE request with data body
// to the URL built with subDomain and path.
func (c Client) Delete(subDomain string, path string, data interface{}) (body []byte, header http.Header, err error) {
	body, header, err = c.process("DELETE", subDomain, path, data)
	return
}

// EmptyTokenError occurs when request is sent without token.
type EmptyTokenError struct{}

func (err EmptyTokenError) Error() (msg string) {
	msg = fmt.Sprintf("empty token: publish personal access token at https://qiita.com/settings/applications, then set environment variable as %s", envAccessToken)
	return
}

// WrongTokenError occurs when the sent token is invalid.
type WrongTokenError struct{}

func (err WrongTokenError) Error() (msg string) {
	msg = fmt.Sprintf("wrong token: publish personal access token at https://qiita.com/settings/applications, then set environment variable as %s", envAccessToken)
	return
}

// ResponseError occurs when the response status is failed
// and the body is JSON.
type ResponseError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (err ResponseError) Error() (msg string) {
	msg = err.Message
	return
}

// StatusError occurs when the response status is failed
// and the body isn't JSON.
type StatusError struct {
	Code    int
	Message string
}

func (err StatusError) Error() (msg string) {
	msg = err.Message
	return
}
