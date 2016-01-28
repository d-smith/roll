package main

import (
	"fmt"
	"github.com/xtraclabs/roll/jwtsample"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		println("usage: go rollTokenFromAuth0Token <auth0Token>")
		return
	}
	jwtResponse := jwtsample.TradeTokenForToken(os.Args[1], "http://localhost:3000/oauth2/token")
	fmt.Println("\n", jwtResponse)
}
