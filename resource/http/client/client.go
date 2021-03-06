package client

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/alileza/tomato/config"
)

type response struct {
	Code int
	Body []byte
}

type Client struct {
	httpClient   *http.Client
	baseURL      string
	lastResponse *response
}

func New(cfg *config.Resource) (*Client, error) {
	params := cfg.Params

	client := &Client{new(http.Client), "", nil}
	for key, val := range params {
		switch key {
		case "base_url":
			client.baseURL = val
		case "timeout":
			timeout, err := time.ParseDuration(val)
			if err != nil {
				return nil, errors.New("timeout: get http client, invalid params value : " + err.Error())
			}
			client.httpClient.Timeout = timeout
		default:
			return nil, errors.New(key + ": invalid params")
		}
	}
	return client, nil
}

func (c *Client) Ready() error {
	resp, err := c.httpClient.Get(c.baseURL)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusServiceUnavailable {
		return errors.New("http/client: server unavailable > " + c.baseURL)
	}
	return nil
}

func (c *Client) Reset() error {
	c.lastResponse = nil
	return nil
}

func (c *Client) Response() (int, []byte, error) {
	if c.lastResponse == nil {
		return 0, nil, errors.New("no request has been sent, please send request before checking response")
	}
	return c.lastResponse.Code, c.lastResponse.Body, nil
}

func (c *Client) Request(method, path string, body []byte) error {
	req, err := http.NewRequest(method, path, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	if c.baseURL != "" {
		baseURL, err := url.Parse(c.baseURL)
		if err != nil {
			return err
		}
		req.URL.Scheme = baseURL.Scheme
		req.URL.Host = baseURL.Host
	}

	if req.Method != http.MethodGet {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	c.lastResponse = &response{resp.StatusCode, body}
	return nil
}
