package main

import (
	"fmt"
	"flag"
	"github.com/xtraclabs/roll/rollsvcs"
)

func main() {

	var port = flag.Int("port", -1, "Port to listen on")
	flag.Parse()
	if *port == -1 {
		fmt.Println("Must specify a -port argument")
		return
	}

	coreConfig := rollsvcs.DefaultConfig()
	rollsvcs.RunRoll(*port,coreConfig)
}
