package doh

import "net/http"

// Option represents configuration options for doh.Client.
type Option interface {
	apply(c *Client)
}

type httpClientOption struct {
	client *http.Client
}

func (o *httpClientOption) apply(c *Client) {
	c.client = o.client
}

// WithHTTPClient is a configuration option that overrides default http.Client instance used by the doh.Client.
func WithHTTPClient(c *http.Client) Option {
	return &httpClientOption{
		client: c,
	}
}
