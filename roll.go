package main

import (
	rollhttp "github.com/xtraclabs/roll/http"
	"github.com/xtraclabs/roll/roll"
	"log"
	"net/http"
)

func main() {
	var coreConfig = roll.CoreConfig{}
	core := roll.NewCore(&coreConfig)

	log.Println("Listening on port 12345")
	http.ListenAndServe(":12345", rollhttp.Handler(core))
}
