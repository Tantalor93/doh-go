[![Go Report Card](https://goreportcard.com/badge/github.com/Tantalor93/doh-go)](https://goreportcard.com/report/github.com/Tantalor93/doh-go)
[![Tantalor93](https://circleci.com/gh/Tantalor93/doh-go/tree/main.svg?style=svg)](https://circleci.com/gh/Tantalor93/doh-go?branch=main)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/Tantalor93/doh-go/blob/main/LICENSE)
[![codecov](https://codecov.io/gh/Tantalor93/doh-go/branch/main/graph/badge.svg?token=MC6PK2OLMK)](https://codecov.io/gh/Tantalor93/doh-go)
[![](https://godoc.org/github.com/Tantalor93/doh-go/doh?status.svg)](https://godoc.org/github.com/tantalor93/doh-go/doh)

# doh-go
DoH client written in Golang with minimal dependencies, built on top of https://github.com/miekg/dns
and standard http client (net/http) and based on [DoH RFC](https://datatracker.ietf.org/doc/html/rfc8484#section-4.1)

## Usage in your project
add dependency
```
go get github.com/tantalor93/doh-go
```

## Examples
```
// create client
c := doh.NewClient(nil)

// prepare payload
msg := dns.Msg{}
msg.SetQuestion("google.com.", dns.TypeA)

// send DNS query to Cloudflare Server over DoH using POST method
r, err := c.SendViaPost(context.Background(), "https://1.1.1.1/dns-query", &msg)
if err != nil {
    panic(err)
}

// do something with response
fmt.Println(dns.RcodeToString[r.Rcode])

// send DNS query to Cloudflare Server over DoH using GET method
r, err = c.SendViaGet(context.Background(), "https://1.1.1.1/dns-query", &msg)
if err != nil {
    panic(err)
}

// do something with response
fmt.Println(dns.RcodeToString[r.Rcode])
```