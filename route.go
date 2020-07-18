package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Method  string
	Pattern string
	Handler http.HandlerFunc
	Kind    string
	// Middleware mux.MiddlewareFunc
}

var routes []Route

func init() {
	register("GET", "/movies", AllMovies, "needMiddle")
	register("GET", "/movies/{id}", FindMovie, "")
	register("POST", "/movies", CreateMovie, "")
	register("PUT", "/movies", UpdateMovie, "")
	register("DELETE", "/movies", DeleteMovie, "")
}

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	for _, route := range routes {
		r.HandleFunc(route.Pattern, route.Handler).Methods(route.Method)

		switch route.Kind {
		case "needMiddle":
			r.Use(MiddlewareOne)
		default:
		}
	}
	return r
}

func register(method, pattern string, handler http.HandlerFunc, kind string) {
	routes = append(routes, Route{method, pattern, handler, kind})
}
