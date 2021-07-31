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

// send to Cloudflare Server over DoH
r, err := c.PostSend(context.Background(), "https://1.1.1.1", &msg)

if err != nil {
    panic(err)
}

// do something with response
fmt.Println(dns.RcodeToString[r.Rcode])
```