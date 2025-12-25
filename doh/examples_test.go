package doh_test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/miekg/dns"
	"github.com/tantalor93/doh-go/doh"
	"golang.org/x/net/http2"
)

func ExampleClient_SendViaPost() {
	// create client with default settings resolving via CloudFlare DoH Server
	c := doh.NewClient("https://1.1.1.1/dns-query")

	// prepare payload
	msg := dns.Msg{}
	msg.SetQuestion("google.com.", dns.TypeA)

	// send DNS query using HTTP POST method
	r, err := c.SendViaPost(context.Background(), &msg)
	if err != nil {
		panic(err)
	}

	// do something with response
	fmt.Println(dns.RcodeToString[r.Rcode])
	// Output: NOERROR
}

func ExampleClient_SendViaGet() {
	// create client with default settings resolving via CloudFlare DoH Server
	c := doh.NewClient("https://1.1.1.1/dns-query")

	// prepare payload
	msg := dns.Msg{}
	msg.SetQuestion("google.com.", dns.TypeA)

	// send DNS query using HTTP POST method
	r, err := c.SendViaGet(context.Background(), &msg)
	if err != nil {
		panic(err)
	}

	// do something with response
	fmt.Println(dns.RcodeToString[r.Rcode])
	// Output: NOERROR
}

func ExampleClient_http2() {
	// create client with default settings resolving via CloudFlare DoH Server
	httpClient := http.Client{Transport: &http2.Transport{}}
	c := doh.NewClient("https://1.1.1.1/dns-query", doh.WithHTTPClient(&httpClient))

	// prepare payload
	msg := dns.Msg{}
	msg.SetQuestion("google.com.", dns.TypeA)

	// send DNS query using HTTP POST method
	r, err := c.SendViaGet(context.Background(), &msg)
	if err != nil {
		panic(err)
	}

	// do something with response
	fmt.Println(dns.RcodeToString[r.Rcode])
	// Output: NOERROR
}
