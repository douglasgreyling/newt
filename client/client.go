package client

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
)

type Options struct {
	Headers             Headers
	SkipSSLVerification bool
}

type Headers map[string]string

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

func (c *Client) Get(endpoint string, params ...url.Values) (*http.Response, error) {
    var queryParams url.Values

    if len(params) > 0 {
        queryParams = params[0]
    }

	req, err := c.newRequest(http.MethodGet, endpoint, queryParams, nil)
	if err != nil {
		return nil, err
	}

	return c.do(req)
}

func (c *Client) newRequest(method, endpoint string, params url.Values, body io.Reader) (*http.Request, error) {
	rel, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	fullURL := c.BaseURL.ResolveReference(rel)

    if params != nil {
        fullURL.RawQuery = params.Encode()
    }

	req, err := http.NewRequest(method, fullURL.String(), body)
	if err != nil {
		return nil, err
	}

	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	return c.HTTPClient.Do(req)
}
