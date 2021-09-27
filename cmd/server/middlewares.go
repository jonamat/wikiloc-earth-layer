package main

import (
	"log"
	"net/http"
)

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request. Method: %s. URI: %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func webClientHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "earth.google.com")
		next.ServeHTTP(w, r)
	})
}
