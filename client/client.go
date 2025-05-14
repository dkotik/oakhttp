/*
Package client provides a more secure standard HTTP client.

Higher security is achieved by setting short reasonable timeouts. Default HTTP client [leaks go routines][read-response] if the `response.Body` is not fully read and closed on **every** request or if the server hangs the connection or responds slowly.

## Default HTTP Client

> Making things worse is the fact that a bare http.Client will use a default http.Transport called http.DefaultTransport, which is another global value that behaves the same way. So it is not simply enough to replace http.DefaultClient with &http.Client{}. (https://pkg.go.dev/github.com/hashicorp/go-cleanhttp#section-readme)

TODO: Use Hashicorp's `go-cleanhttp` package to get true copies of http.Client{}.
TODO: Include retriable client from here: https://github.com/projectdiscovery/retryablehttp-go

[read-response]: https://manishrjain.com/must-close-golang-http-response
*/
package client

import (
	"errors"
	"fmt"
	"net"
	"net/http"
)

func New(withOptions ...Option) (*http.Client, error) {
	o := &options{}
	var err error
	for _, option := range append(withOptions, WithDefaultOptions()) {
		if err = option(o); err != nil {
			return nil, fmt.Errorf("cannot initialize a secure HTTP client: %w", err)
		}
	}

	if o.MaxIdleConnsPerHost > o.MaxConnsPerHost {
		return nil, errors.New("cannot have more MaxIdleConnsPerHost than MaxConnsPerHost")
	}

	if o.TLSHandshakeTimeout+o.ResponseHeaderTimeout > o.Timeout {
		return nil, errors.New("sum of header timeouts must not exceed the total Timeout")
	}

	return &http.Client{
		Timeout: o.Timeout,
		Transport: &http.Transport{
			MaxConnsPerHost:     o.MaxConnsPerHost,
			MaxIdleConnsPerHost: o.MaxIdleConnsPerHost,
			DialContext: (&net.Dialer{
				Timeout:   o.Timeout,
				KeepAlive: o.KeepAlive,
			}).DialContext,
			TLSHandshakeTimeout:   o.TLSHandshakeTimeout,
			ResponseHeaderTimeout: o.ResponseHeaderTimeout,
			ExpectContinueTimeout: o.ExpectContinueTimeout,
		}}, nil
}
