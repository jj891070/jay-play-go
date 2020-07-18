package main

import (
	"fmt"
	"net/http"
)

func AllMovies(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented !")
}

func FindMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented !")
}

func CreateMovie(w http.ResponseWriter, t *http.Request) {
	fmt.Fprintln(w, "not implemented !")
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented !")
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented !")
}
