package doh

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const(
	existingDomain = "google.com."
	notExistingDomain = "nxdomain.cz."
)

func Test_SendViaPost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bd, err := ioutil.ReadAll(r.Body)

		msg := dns.Msg{}
		err = msg.Unpack(bd)
		require.NoError(t, err, "error unpacking request body")
		require.Len(t, msg.Question, 1, "single question expected")

		resp := msg
		switch msg.Question[0].Name {
		case notExistingDomain:
			resp.Rcode = dns.RcodeNameError
		case existingDomain:
			resp.Rcode = dns.RcodeSuccess
		default:
			require.FailNow(t, "unexpected question name")
		}

		pack, err := resp.Pack()
		require.NoError(t, err, "error packing response")

		_, err = w.Write(pack)
		require.NoError(t, err, "error writing response")
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
		require.NotEmpty(t, dnsQryParam, "expected dns query param not found")

		bd, err := base64.StdEncoding.DecodeString(dnsQryParam)
		require.NoError(t, err, "error decoding query param DNS")

		msg := dns.Msg{}
		err = msg.Unpack(bd)
		require.NoError(t, err, "error unpacking request body")
		require.Len(t, msg.Question, 1, "single question expected")

		resp := msg
		switch msg.Question[0].Name {
		case notExistingDomain:
			resp.Rcode = dns.RcodeNameError
		case existingDomain:
			resp.Rcode = dns.RcodeSuccess
		default:
			require.FailNow(t, "unexpected question name")
		}

		pack, err := resp.Pack()
		require.NoError(t, err, "error packing response")

		_, err = w.Write(pack)
		require.NoError(t, err, "error writing response")
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
