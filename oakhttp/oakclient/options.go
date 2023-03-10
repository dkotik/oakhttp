package oakclient

import (
	"errors"
	"time"
)

type options struct {
	MaxConnsPerHost       int
	MaxIdleConnsPerHost   int
	Timeout               time.Duration
	KeepAlive             time.Duration
	TLSHandshakeTimeout   time.Duration
	ResponseHeaderTimeout time.Duration
	ExpectContinueTimeout time.Duration
}

type Option func(*options) error

// WithTimeout sets [http.Client.Timeout] and [net.Dialer.Timeout].
func WithTimeout(d time.Duration) Option {
	return func(o *options) error {
		if o.Timeout != 0 {
			return errors.New("Timeout is already set")
		}
		if o.Timeout < time.Second {
			return errors.New("Timeout must be greater than one second")
		}
		if o.Timeout > time.Minute*5 {
			return errors.New("Timeout must be less than five minutes")
		}
		o.Timeout = d
		return nil
	}
}

// WithKeepAlive sets [net.Dialer.KeepAlive].
func WithKeepAlive(d time.Duration) Option {
	return func(o *options) error {
		if o.KeepAlive != 0 {
			return errors.New("KeepAlive is already set")
		}
		if o.KeepAlive < time.Second {
			return errors.New("KeepAlive must be greater than one second")
		}
		if o.KeepAlive > time.Hour*24 {
			return errors.New("KeepAlive must be less than a day")
		}
		o.KeepAlive = d
		return nil
	}
}

// WithTLSHandshakeTimeout sets [http.Transport.TLSHandshakeTimeout].
func WithTLSHandshakeTimeout(d time.Duration) Option {
	return func(o *options) error {
		if o.TLSHandshakeTimeout != 0 {
			return errors.New("TLSHandshakeTimeout is already set")
		}
		if o.TLSHandshakeTimeout < time.Second {
			return errors.New("TLSHandshakeTimeout must be greater than one second")
		}
		if o.TLSHandshakeTimeout > time.Second*10 {
			return errors.New("TLSHandshakeTimeout must be less than ten seconds")
		}
		o.TLSHandshakeTimeout = d
		return nil
	}
}

// WithResponseHeaderTimeout sets [http.Transport.ResponseHeaderTimeout].
func WithResponseHeaderTimeout(d time.Duration) Option {
	return func(o *options) error {
		if o.ResponseHeaderTimeout != 0 {
			return errors.New("ResponseHeaderTimeout is already set")
		}
		if o.ResponseHeaderTimeout < time.Second {
			return errors.New("ResponseHeaderTimeout must be greater than one second")
		}
		if o.ResponseHeaderTimeout > time.Second*10 {
			return errors.New("ResponseHeaderTimeout must be less ten seconds")
		}
		o.ResponseHeaderTimeout = d
		return nil
	}
}

// WithExpectContinueTimeout sets [http.Transport.ExpectContinueTimeout].
func WithExpectContinueTimeout(d time.Duration) Option {
	return func(o *options) error {
		if o.ExpectContinueTimeout != 0 {
			return errors.New("ExpectContinueTimeout is already set")
		}
		if o.ExpectContinueTimeout < time.Second {
			return errors.New("ExpectContinueTimeout must be greater than one second")
		}
		if o.ExpectContinueTimeout > time.Second*10 {
			return errors.New("ExpectContinueTimeout must be less than ten seconds")
		}
		o.ExpectContinueTimeout = d
		return nil
	}
}

// WithMaxConnsPerHost sets [http.Transport.MaxConnsPerHost].
func WithMaxConnsPerHost(n int) Option {
	return func(o *options) error {
		if o.MaxConnsPerHost != 0 {
			return errors.New("MaxConnsPerHost is already set")
		}
		if o.MaxConnsPerHost < 1 {
			return errors.New("MaxConnsPerHost must be greater than 0")
		}
		if o.MaxConnsPerHost > 1024 {
			return errors.New("MaxConnsPerHost greater than 1024 is unreasonable")
		}
		o.MaxConnsPerHost = n
		return nil
	}
}

// WithMaxIdleConnsPerHost sets [http.Transport.MaxIdleConnsPerHost].
func WithMaxIdleConnsPerHost(n int) Option {
	return func(o *options) error {
		if o.MaxIdleConnsPerHost != 0 {
			return errors.New("MaxIdleConnsPerHost is already set")
		}
		if o.MaxIdleConnsPerHost < 1 {
			return errors.New("MaxIdleConnsPerHost must be greater than 0")
		}
		if o.MaxIdleConnsPerHost > 1024 {
			return errors.New("MaxIdleConnsPerHost greater than 1024 is unreasonable")
		}
		o.MaxIdleConnsPerHost = n
		return nil
	}
}
