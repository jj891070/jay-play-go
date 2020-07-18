package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"html/template"
)

var clientID string
var clientSecret string
var callbackURL string
var token string

func main() {
	http.HandleFunc("/callback", callbackHandler)
	http.HandleFunc("/notify", notifyHandler)
	http.HandleFunc("/auth", authHandler)
	clientID = os.Getenv("ClientID")
	clientSecret = os.Getenv("ClientSecret")
	callbackURL = os.Getenv("CallbackURL")
	port := os.Getenv("PORT")
	fmt.Printf("ENV port:%s, cid:%s csecret:%s\n", port, clientID, clientSecret)
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func notifyHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // Populates request.Form
	msg := r.Form.Get("msg")
	fmt.Printf("Get msg=%s\n", msg)

	data := url.Values{}
	data.Add("message", msg)

	byt, err := apiCall("POST", apiNotify, data, token)
	fmt.Println("ret:", string(byt), " err:", err)

	res := newTokenResponse(byt)
	fmt.Println("result:", res)
	token = res.AccessToken
	w.Write(byt)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // Populates request.Form
	code := r.Form.Get("code")
	state := r.Form.Get("state")
	fmt.Printf("Get code=%s, state=%s \n", code, state)

	data := url.Values{}
	data.Add("grant_type", "authorization_code")
	data.Add("code", code)
	data.Add("redirect_uri", callbackURL)
	data.Add("client_id", clientID)
	data.Add("client_secret", clientSecret)

	byt, err := apiCall("POST", apiToken, data, "")
	fmt.Println("ret:", string(byt), " err:", err)

	res := newTokenResponse(byt)
	fmt.Println("result:", res)
	token = res.AccessToken
	w.Write(byt)
}
func authHandler(w http.ResponseWriter, r *http.Request) {
	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}
	t, err := template.New("webpage").Parse(authTmpl)
	check(err)
	noItems := struct {
		ClientID    string
		CallbackURL string
	}{
		ClientID:    clientID,
		CallbackURL: callbackURL,
	}

	err = t.Execute(w, noItems)
	check(err)
}
