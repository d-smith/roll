package main

import (
	rollhttp "github.com/xtraclabs/roll/http"
	"github.com/xtraclabs/roll/roll"
	"log"
	"net/http"
	"github.com/xtraclabs/roll/repos"
	"flag"
	"fmt"
)

func main() {

	var port = flag.Int("port", -1, "Port to listen on")
	flag.Parse()
	if *port == -1 {
		fmt.Println("Must specify a -port argument")
		return
	}

	var coreConfig = roll.CoreConfig{
		DeveloperRepo: repos.NewDynamoDevRepo(),
		ApplicationRepo: repos.NewDynamoAppRepo(),
	}


	core := roll.NewCore(&coreConfig)

	log.Println("Listening on port ",*port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), rollhttp.Handler(core))
}
