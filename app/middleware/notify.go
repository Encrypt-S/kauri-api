package middleware

import (
	"log"
	"net/http"
)

// Notify function adapts an http.Handler to write out “before” and “after” strings
// allowing the original http.Handler `h` to do whatever it was already going to do in between.
// it returns an Adapter, which is just a function that takes and returns an http.Handler.
func Notify() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("before")
			defer log.Println("after")
			h.ServeHTTP(w, r)
		})
	}
}

// usage example for Notify ()
// logger := log.New(os.Stdout, "server: ", log.Lshortfile)
// http.Handle("/", Notify(logger)(indexHandler))
