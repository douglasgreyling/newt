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
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	client := &Client{
		BaseURL:             parsedURL,
		Headers:             options[0].Headers,
		SkipSSLVerification: options[0].SkipSSLVerification,
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

func (c *Client) Get(endpoint string) (*http.Response, error) {
	req, err := c.newRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return c.do(req)
}

func (c *Client) newRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	rel, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	fullURL := c.BaseURL.ResolveReference(rel)

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
