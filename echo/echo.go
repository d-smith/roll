package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

func echoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body.Close()
			w.Write(body)
			w.Write([]byte("\n"))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
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
	mux.Handle("/echo", echoHandler())
	http.ListenAndServe(fmt.Sprintf(":%d", *port), mux)
}
