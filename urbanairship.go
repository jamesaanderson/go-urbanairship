package urbanairship

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	apiVersion = "3"
	baseURL    = "https://go.urbanairship.com"
	mimeType   = "application/vnd.urbanairship+json; version=" + apiVersion
)

type Client struct {
	client   *http.Client
	BaseURL  *url.URL
	MimeType string
	Key      string
	Secret   string
}

func NewClient(httpClient *http.Client, key, secret string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(baseURL)

	return &Client{
		client:   httpClient,
		BaseURL:  baseURL,
		MimeType: mimeType,
		Key:      key,
		Secret:   secret,
	}
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", c.MimeType)
	req.SetBasicAuth(c.Key, c.Secret)
	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return resp, err
}

type ErrorResponse struct {
	Response  *http.Response
	Ok        bool         `json:"ok"`
	ErrorStr  string       `json:"error"`
	ErrorCode int          `json:"error_code"`
	Details   ErrorDetails `json:"details"`
}

type ErrorDetails struct {
	Error    string        `json:"error"`
	Path     string        `json:"path"`
	Location ErrorLocation `json:"location"`
}

type ErrorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

func (r ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.ErrorStr)
}

func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}
