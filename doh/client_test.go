package doh

import (
	"context"
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

func Test_SendViaPost(t *testing.T) {
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
			args:      args{server: "https://1.1.1.1/dns-query", msg: question("google.com.")},
			wantRcode: dns.RcodeSuccess,
		},
		{
			name:      "NXDOMAIN DNS resolution ",
			args:      args{server: "https://1.1.1.1/dns-query", msg: question("nxdomain.cz.")},
			wantRcode: dns.RcodeNameError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(nil)

			got, err := client.SendViaPost(context.Background(), tt.args.server, tt.args.msg)

			if tt.wantErr {
				assert.Error(t, err, "SendViaPost() error")
			} else {
				assert.NotNil(t, got, "SendViaPost() response")
				assert.Equal(t, tt.wantRcode, got.Rcode, "SendViaPost() rcode")
			}
		})
	}
}

func Test_SendViaGet(t *testing.T) {
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
			args:      args{server: "https://1.1.1.1/dns-query", msg: question("google.com.")},
			wantRcode: dns.RcodeSuccess,
		},
		{
			name:      "NXDOMAIN DNS resolution ",
			args:      args{server: "https://1.1.1.1/dns-query", msg: question("nxdomain.cz.")},
			wantRcode: dns.RcodeNameError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(nil)

			got, err := client.SendViaGet(context.Background(), tt.args.server, tt.args.msg)

			if tt.wantErr {
				assert.Error(t, err, "SendViaGet() error")
			} else {
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
