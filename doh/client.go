package doh

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/miekg/dns"
)

// Client encapsulates and provides logic for querying DNS servers over DoH
type Client struct {
	c *http.Client
}

// NewClient creates new Client instance with standard net/http client. If nil, default http.Client is used.
func NewClient(c *http.Client) *Client {
	if c == nil {
		c = &http.Client{}
	}
	return &Client{c}
}

// SendViaPost sends DNS message to the given DNS server over DoH using POST, see https://datatracker.ietf.org/doc/html/rfc8484#section-4.1
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

// SendViaGet sends DNS message to the given DNS server over DoH using GET, see https://datatracker.ietf.org/doc/html/rfc8484#section-4.1
func (dc *Client) SendViaGet(ctx context.Context, server string, msg *dns.Msg) (*dns.Msg, error) {
	pack, err := msg.Pack()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprint(server, "?dns=", base64.URLEncoding.EncodeToString(pack))
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request = request.WithContext(ctx)
	request.Header.Set("Accept", "application/dns-message")

	return dc.send(request)
}

func (dc *Client) send(r *http.Request) (*dns.Msg, error) {
	resp, err := dc.c.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected HTTP status")
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
