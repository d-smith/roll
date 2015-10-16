package main

import (
	"flag"
	"fmt"
	rollhttp "github.com/xtraclabs/roll/http"
	"github.com/xtraclabs/roll/repos"
	"github.com/xtraclabs/roll/roll"
	"log"
	"net/http"
)

func main() {

	var port = flag.Int("port", -1, "Port to listen on")
	flag.Parse()
	if *port == -1 {
		fmt.Println("Must specify a -port argument")
		return
	}

	var coreConfig = roll.CoreConfig{
		DeveloperRepo:   repos.NewDynamoDevRepo(),
		ApplicationRepo: repos.NewDynamoAppRepo(),
		SecretsRepo:     repos.NewVaultSecretsRepo(),
		IdGenerator: new(roll.UUIDIdGenerator),
	}

	core := roll.NewCore(&coreConfig)

	log.Println("Listening on port ", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), rollhttp.Handler(core))
}
