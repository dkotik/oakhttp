package client

import (
	"fmt"
	"time"
)

// WithDefaultOptions sets any unset [Option] to its default value. The default values are optimized for a typical secure remote procedure calls.
func WithDefaultOptions() Option {
	return func(o *options) (err error) {
		defer func() {
			if err != nil {
				err = fmt.Errorf("could not set default setting: %w", err)
			}
		}()

		if o.MaxConnsPerHost == 0 {
			if err = WithConnectionLimit(20)(o); err != nil {
				return err
			}
		}
		if o.MaxIdleConnsPerHost == 0 {
			if err = WithIdleConnectionLimitPerHost(
				int(float32(o.MaxConnsPerHost)*.2) + 1,
			)(o); err != nil {
				return err
			}
		}

		if o.Timeout == 0 {
			if err = WithTimeout(time.Second * 5)(o); err != nil {
				return err
			}
		}
		if o.KeepAlive == 0 {
			if err = WithKeepAliveTimeout(time.Second * 60)(o); err != nil {
				return err
			}
		}
		if o.TLSHandshakeTimeout == 0 {
			if err = WithTLSHandshakeTimeout(time.Second * 2)(o); err != nil {
				return err
			}
		}
		if o.ResponseHeaderTimeout == 0 {
			if err = WithResponseHeaderTimeout(time.Second * 2)(o); err != nil {
				return err
			}
		}
		if o.ExpectContinueTimeout == 0 {
			if err = WithExpectContinueTimeout(time.Second * 1)(o); err != nil {
				return err
			}
		}
		return nil
	}
}
