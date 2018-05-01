package middleware

import (
	"net/http"
)

// Adapter type function is a wrapper that takes in and returns an http.Handler
type Adapter func(http.Handler) http.Handler

// existing http.Handler is passed in and the Adapter will adapt it
// and return a new (wrapped) http.Handler to use in its place
// To make the adapters run in the order in which they are specified
// you could reverse through them in the Adapt function

// in essence, this middleware allows us to run code before
// and/or after our handler code in a HTTP request lifecycle
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}
