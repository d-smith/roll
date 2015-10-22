package integrationtest

import (
	. "github.com/lsegal/gucumber"
	"log"
	"github.com/xtraclabs/roll/internal/runutils"
	"syscall"
	"time"
)

func init() {
	var shutdownDone chan bool

	Given(`^Some tests to run$`, func() {
	})

	Then(`^The test process and supporting containers are started before and stopped after$`, func() {
	})

	GlobalContext.BeforeAll(func() {
		log.Println("starting test containers")
		shutdownDone = runutils.RunVaultAndRoll()
		//TODO - need to signal when the container is ready... for now sleep...
		time.Sleep(1 * time.Second)
	})

	GlobalContext.AfterAll(func() {
		log.Println("stop and remove test containers")
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		<-shutdownDone
		//Send signal
	})


}