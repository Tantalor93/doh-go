package doh

import "fmt"

// UnexpectedServerHTTPStatusError error indicating that DoH server responded with bad HTTP status code.
type UnexpectedServerHTTPStatusError struct {
	code int
}

func (u UnexpectedServerHTTPStatusError) Error() string {
	return fmt.Sprintf("unexpected upstream server response HTTP status: %d", u.code)
}

// HTTPStatus HTTP status code returned by the DoH Server.
func (u UnexpectedServerHTTPStatusError) HTTPStatus() int {
	return u.code
}
