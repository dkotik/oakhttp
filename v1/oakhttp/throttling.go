package oathttp

import "net/http"

type HTTPThrottler func(r *http.Request) error
