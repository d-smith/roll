package main

import (
	"github.com/xtraclabs/roll/internal/runutils"
)

func main() {
	shutdownDone := runutils.RunVaultAndRoll()
	<-shutdownDone
}
