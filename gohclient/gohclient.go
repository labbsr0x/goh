package gohclient

import (
	"bytes"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// API defines an interface for helper methods that encapsulates http requests complexities
type API interface {
	Put(url string, data []byte) (*http.Response, []byte, error)
	Post(url string, data []byte) (*http.Response, []byte, error)
	Get(url string) (*http.Response, []byte, error)
	Delete(url string) (*http.Response, []byte, error)
}

// Default defines a struct that handles with HTTP requests for a bindman webhook client
type Client struct {
	// User agent used when communicating with the API
	UserAgent string
	// Request content type used when communicating with the API
	ContentType string
	Accept      string
	baseURL     *url.URL
	httpClient  *http.Client
}

// New instantiates a default goh client
// If a nil httpClient is provided, http.DefaultClient will be used.
func NewClient(httpClient *http.Client, baseURL string) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if strings.TrimSpace(baseURL) == "" {
		return nil, errors.New("base URL cannot be an empty string")
	}
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		baseURL:    parsedURL,
		httpClient: httpClient,
	}, nil
}

// Put wraps the call to http.NewRequest apis and properly submits a new HTTP POST request
func (c *Client) Put(url string, data []byte) (*http.Response, []byte, error) {
	return c.request(url, "PUT", data)
}

// Post wraps the call to http.NewRequest apis and properly submits a new HTTP POST request
func (c *Client) Post(url string, data []byte) (*http.Response, []byte, error) {
	return c.request(url, "POST", data)
}

// Get wraps the call to http.NewRequest apis and properly submits a new HTTP GET request
func (c *Client) Get(url string) (*http.Response, []byte, error) {
	return c.request(url, "GET", nil)
}

// Delete wraps the call to http.NewRequest apis and properly submits a new HTTP DELETE request
func (c *Client) Delete(url string) (*http.Response, []byte, error) {
	return c.request(url, "DELETE", nil)
}

// request defines a generic method to execute http requests
func (c *Client) request(path, method string, body []byte) (resp *http.Response, data []byte, err error) {
	if c.httpClient == nil {
		err = fmt.Errorf("httpClient cannot be nil. Use gohclient.NewClient to create a correct Client instance")
		return
	}
	if c.baseURL == nil {
		err = fmt.Errorf("baseURL cannot be nil. Use gohclient.NewClient to create a correct Client instance")
		return
	}

	u, err := c.baseURL.Parse(path)
	if err != nil {
		return
	}

	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(body))
	if err != nil {
		logrus.Errorf("HTTP request creation failed. err=%v", err)
		return
	}

	if body != nil && strings.TrimSpace(c.ContentType) != "" {
		req.Header.Set("Content-Type", c.ContentType)
	}
	if strings.TrimSpace(c.Accept) != "" {
		req.Header.Set("Accept", c.Accept)
	}
	if strings.TrimSpace(c.UserAgent) != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	logrus.Debugf("%v request=%v", method, req)

	resp, err = c.httpClient.Do(req)
	if err != nil {
		logrus.Errorf("HTTP  %v request invocation failed. err=%v", method, err)
		return
	}
	defer dClose(resp.Body)
	logrus.Debugf("Response: %v", resp)
	data, err = ioutil.ReadAll(resp.Body)
	logrus.Debugf("Response body: %v", data)
	return
}

func dClose(c io.Closer) {
	if err := c.Close(); err != nil {
		logrus.Errorf("HTTP response body close invocation failed. err=%v", err)
	}
}
