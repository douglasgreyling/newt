package client

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"
)

type Options struct {
	Headers             Headers
	SkipSSLVerification bool
}

type Headers map[string]string

type Params map[string]string

type Client struct {
	BaseURL             *url.URL
	HTTPClient          *http.Client
	Headers             Headers
	SkipSSLVerification bool
}

func NewClient(baseURL string, options ...Options) (*Client, error) {
	var opts Options

	if len(options) > 0 {
		opts = options[0]
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	client := &Client{
		BaseURL:             parsedURL,
		Headers:             opts.Headers,
		SkipSSLVerification: opts.SkipSSLVerification,
	}

	client.HTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: client.SkipSSLVerification,
			},
		},
	}

	return client, nil
}

func (c *Client) Get(endpoint string, params ...Params) (*http.Response, error) {
	var prms Params

	if len(params) > 0 {
		prms = params[0]
	}

	req, err := c.newRequest(http.MethodGet, endpoint, prms)
	if err != nil {
		return nil, err
	}

	return c.do(req)
}

func (c *Client) Post(endpoint string, params Params) (*http.Response, error) {
	req, err := c.newRequest(http.MethodPost, endpoint, params)
	if err != nil {
		return nil, err
	}

	return c.do(req)
}

func (c *Client) newRequest(method, endpoint string, params Params) (*http.Request, error) {
	parsedUrl, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	resolvedUrl := c.BaseURL.ResolveReference(parsedUrl)
	formData := c.parseParams(params).Encode()

	req, err := http.NewRequest(method, resolvedUrl.String(), strings.NewReader(formData))
	if err != nil {
		return nil, err
	}

	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	if c.Headers["Content-Type"] == "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return req, nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	return c.HTTPClient.Do(req)
}

func (c *Client) parseParams(params Params) url.Values {
	body := make(url.Values)

	for key, value := range params {
		body.Add(key, value)
	}

	return body
}
