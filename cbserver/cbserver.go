package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"io/ioutil"
	"os"
)

var templates = template.Must(template.ParseFiles("static/callback.html"))

var azServerEndpoint string
var clientId string
var clientSecret string
var redirectURI string

func init() {
	azServerEndpoint = os.Getenv("AZ_SERVER")
	fmt.Println("AZ_SERVER:", azServerEndpoint)

	clientId = os.Getenv("CLIENT_ID")
	fmt.Println("CLIENT_ID:", clientId)

	clientSecret = os.Getenv("CLIENT_SECRET")
	fmt.Println("CLIENT_SECRET:", clientSecret)

	redirectURI = os.Getenv("REDIRECT_URI")
	fmt.Println("REDIRECT_URI:", redirectURI)
}

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
			doAuthCodeCallback(w,r)
		}
	}
}

func doAuthCodeCallback(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	codes := params["code"]

	resp, err := http.PostForm("http://" + azServerEndpoint + "/oauth2/token",
		url.Values{"grant_type" : {"authorization_code"},
		"code":{codes[0]}, "client_id":{clientId}, "client_secret":{clientSecret},
		"redirect_uri":{redirectURI}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println(fmt.Sprintf("%v", string(body)))

	w.Write(body)

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
