package main

import (
	"log"
	"net/http"
)

// MiddlewareOne 中間件
func MiddlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middlewareOne")
		next.ServeHTTP(w, r)
		log.Println("Executing middlewareOne again")
	})
}
