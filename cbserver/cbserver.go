package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

var templates = template.Must(template.ParseFiles("static/callback.html", "static/echo.html"))

var azServerEndpoint string
var clientID string
var clientSecret string
var redirectURI string

func init() {
	azServerEndpoint = os.Getenv("AZ_SERVER")
	fmt.Println("AZ_SERVER:", azServerEndpoint)

	clientID = os.Getenv("CLIENT_ID")
	fmt.Println("CLIENT_ID:", clientID)

	clientSecret = os.Getenv("CLIENT_SECRET")
	fmt.Println("CLIENT_SECRET:", clientSecret)

	redirectURI = os.Getenv("REDIRECT_URI")
	fmt.Println("REDIRECT_URI:", redirectURI)
}

type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func isTokenCallback(r *http.Request) bool {
	params := r.URL.Query()
	codes := params["code"]
	log.Println("isTokenCallback codes: ", codes)
	return !(len(codes) == 1)
}

func oauthCallbackHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if isTokenCallback(r) {
			err := templates.ExecuteTemplate(w, "callback.html", nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			doAuthCodeCallback(w, r)
		}
	}
}

func doAuthCodeCallback(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	codes := params["code"]

	resp, err := http.PostForm("http://"+azServerEndpoint+"/oauth2/token",
		url.Values{"grant_type": {"authorization_code"},
			"code": {codes[0]}, "client_id": {clientID}, "client_secret": {clientSecret},
			"redirect_uri": {redirectURI}})
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

	var jsonResponse accessTokenResponse
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = templates.ExecuteTemplate(w, "echo.html", jsonResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
