/*
Package oakhttp holds utility methods, adapters, and builders for hardening the most common elements of standard library `net/http` package. It aims to add to add security by default and resistance to misconfiguration where they are insufficient.
*/
package oakhttp

import "net/http"

type Handler func(w http.ResponseWriter, r *http.Request) error

type Middleware func(Handler) Handler
