package doh

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	existingDomain    = "google.com."
	notExistingDomain = "nxdomain.cz."
)

func Test_SendViaPost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bd, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		msg := dns.Msg{}
		err = msg.Unpack(bd)
		if err != nil {
			panic(err)
		}

		resp := msg
		switch msg.Question[0].Name {
		case notExistingDomain:
			resp.Rcode = dns.RcodeNameError
		case existingDomain:
			resp.Rcode = dns.RcodeSuccess
		default:
			panic("unexpected question name")
		}

		pack, err := resp.Pack()
		if err != nil {
			panic(err)
		}

		_, err = w.Write(pack)
		if err != nil {
			panic(err)
		}
	}))
	defer ts.Close()

	type args struct {
		server string
		msg    *dns.Msg
	}
	tests := []struct {
		name      string
		args      args
		wantRcode int
		wantErr   bool
	}{
		{
			name:      "NOERROR DNS resolution",
			args:      args{server: ts.URL, msg: question(existingDomain)},
			wantRcode: dns.RcodeSuccess,
		},
		{
			name:      "NXDOMAIN DNS resolution",
			args:      args{server: ts.URL, msg: question(notExistingDomain)},
			wantRcode: dns.RcodeNameError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(nil)

			got, err := client.SendViaPost(context.Background(), tt.args.server, tt.args.msg)

			if tt.wantErr {
				require.Error(t, err, "SendViaPost() error")
			} else {
				require.NoError(t, err)
				assert.NotNil(t, got, "SendViaPost() response")
				assert.Equal(t, tt.wantRcode, got.Rcode, "SendViaPost() rcode")
			}
		})
	}
}

func Test_SendViaGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		dnsQryParam := query.Get("dns")

		bd, err := base64.RawURLEncoding.DecodeString(dnsQryParam)
		if err != nil {
			panic(err)
		}

		msg := dns.Msg{}
		err = msg.Unpack(bd)
		if err != nil {
			panic(err)
		}

		resp := msg
		switch msg.Question[0].Name {
		case notExistingDomain:
			resp.Rcode = dns.RcodeNameError
		case existingDomain:
			resp.Rcode = dns.RcodeSuccess
		default:
			panic("unexpected question name")
		}

		pack, err := resp.Pack()
		if err != nil {
			panic(err)
		}

		_, err = w.Write(pack)
		if err != nil {
			panic(err)
		}
	}))
	defer ts.Close()

	type args struct {
		server string
		msg    *dns.Msg
	}
	tests := []struct {
		name      string
		args      args
		wantRcode int
		wantErr   bool
	}{
		{
			name:      "NOERROR DNS resolution",
			args:      args{server: ts.URL, msg: question(existingDomain)},
			wantRcode: dns.RcodeSuccess,
		},
		{
			name:      "NXDOMAIN DNS resolution",
			args:      args{server: ts.URL, msg: question(notExistingDomain)},
			wantRcode: dns.RcodeNameError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(nil)

			got, err := client.SendViaGet(context.Background(), tt.args.server, tt.args.msg)

			if tt.wantErr {
				require.Error(t, err, "SendViaGet() error")
			} else {
				require.NoError(t, err)
				assert.NotNil(t, got, "SendViaGet() response")
				assert.Equal(t, tt.wantRcode, got.Rcode, "SendViaGet() rcode")
			}
		})
	}
}

func question(fqdn string) *dns.Msg {
	q := dns.Msg{}
	return q.SetQuestion(fqdn, dns.TypeA)
}
