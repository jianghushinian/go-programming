//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"
	"net/http"
)

func WithLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next(w, r)
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func main() {
	http.Handle("/", WithLogging(HelloHandler))
	http.ListenAndServe(":8080", nil)
}
