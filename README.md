[![Release](https://img.shields.io/github/release/Tantalor93/doh-go/all.svg)](https://github.com/tantalor93/doh-go/releases)
[![Go version](https://img.shields.io/github/go-mod/go-version/Tantalor93/doh-go)](https://github.com/Tantalor93/doh-go/blob/master/go.mod#L3)
[![](https://godoc.org/github.com/Tantalor93/doh-go/doh?status.svg)](https://godoc.org/github.com/tantalor93/doh-go/doh)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/Tantalor93/doh-go/blob/main/LICENSE)
[![Tantalor93](https://circleci.com/gh/Tantalor93/doh-go/tree/main.svg?style=svg)](https://circleci.com/gh/Tantalor93/doh-go?branch=main)
[![lint](https://github.com/Tantalor93/doh-go/actions/workflows/lint.yml/badge.svg?branch=main)](https://github.com/Tantalor93/doh-go/actions/workflows/lint.yml)
[![codecov](https://codecov.io/gh/Tantalor93/doh-go/branch/main/graph/badge.svg?token=MC6PK2OLMK)](https://codecov.io/gh/Tantalor93/doh-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/Tantalor93/doh-go)](https://goreportcard.com/report/github.com/Tantalor93/doh-go)

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
// create client with default settings
c := doh.NewClient("https://1.1.1.1/dns-query")

// prepare payload
msg := dns.Msg{}
msg.SetQuestion("google.com.", dns.TypeA)

// send DNS query to Cloudflare Server over DoH using POST method
r, err := c.SendViaPost(context.Background(), &msg)
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
