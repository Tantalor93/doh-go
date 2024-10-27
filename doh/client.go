package doh

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/miekg/dns"
)

// Client encapsulates and provides logic for querying DNS servers over DoH.
type Client struct {
	client *http.Client
}

// NewClient creates new Client instance with standard net/http client.
func NewClient(opts ...Option) *Client {
	client := &Client{
		client: &http.Client{},
	}
	for _, opt := range opts {
		opt.apply(client)
	}
	return client
}

// SendViaPost sends DNS message to the given DNS server over DoH using POST method, see https://datatracker.ietf.org/doc/html/rfc8484#section-4.1
func (dc *Client) SendViaPost(ctx context.Context, server string, msg *dns.Msg) (*dns.Msg, error) {
	pack, err := msg.Pack()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", server, bytes.NewReader(pack))
	if err != nil {
		return nil, err
	}
	request = request.WithContext(ctx)
	request.Header.Set("Accept", "application/dns-message")
	request.Header.Set("content-type", "application/dns-message")

	return dc.send(request)
}

// SendViaGet sends DNS message to the given DNS server over DoH using GET method, see https://datatracker.ietf.org/doc/html/rfc8484#section-4.1
func (dc *Client) SendViaGet(ctx context.Context, server string, msg *dns.Msg) (*dns.Msg, error) {
	pack, err := msg.Pack()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprint(server, "?dns=", base64.RawURLEncoding.EncodeToString(pack))
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request = request.WithContext(ctx)
	request.Header.Set("Accept", "application/dns-message")

	return dc.send(request)
}

func (dc *Client) send(r *http.Request) (*dns.Msg, error) {
	resp, err := dc.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, UnexpectedServerHTTPStatusError{code: resp.StatusCode}
	}

	buffer := bytes.Buffer{}
	_, err = buffer.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	res := dns.Msg{}
	err = res.Unpack(buffer.Bytes())
	if err != nil {
		return nil, err
	}
	return &res, nil
}
