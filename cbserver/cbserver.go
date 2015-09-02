package main
import (
	"net/http"
	"flag"
	"fmt"
	"html/template"
)

var templates = template.Must(template.ParseFiles("static/callback.html"))

func oauthCallbackHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := templates.ExecuteTemplate(w, "callback.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
	http.ListenAndServe(fmt.Sprintf(":%d", *port), mux)
}
