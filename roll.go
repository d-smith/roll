package main

import (
	rollhttp "github.com/xtraclabs/roll/http"
	"github.com/xtraclabs/roll/roll"
	"log"
	"net/http"
)

func main() {
	core := roll.NewCore()

	log.Println("Listening on port 12345")
	http.ListenAndServe(":12345", rollhttp.Handler(core))
}
