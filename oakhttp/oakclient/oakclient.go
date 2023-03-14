/*
Package oakclient provides a secure HTTP client.
*/
package oakclient

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
