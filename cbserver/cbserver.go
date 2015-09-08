package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseFiles("static/callback.html"))

func isTokenCallback(r *http.Request) bool {
	params := r.URL.Query()
	codes := params["code"]
	log.Println("isTokenCallback codes: ", codes)
	return ! (len(codes) == 1)
}

func oauthCallbackHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if isTokenCallback(r) {
			err := templates.ExecuteTemplate(w, "callback.html", nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Not implemented", http.StatusInternalServerError)
		}
	}
}

func loginCallbackHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func main() {

	var port = flag.Int("port", -1, "Port to listen on")
	flag.Parse()
	if *port == -1 {
		fmt.Println("Must specify a -port argument")
		return
	}

	mux := http.NewServeMux()
	mux.Handle("/oauth2_callback", oauthCallbackHandler())
	mux.Handle("/XtracWeb/services/Login", loginCallbackHandler())
	http.ListenAndServe(fmt.Sprintf(":%d", *port), mux)
}
