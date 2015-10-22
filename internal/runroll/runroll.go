package main
import (
	"github.com/xtraclabs/roll/internal/dockerutil"
	"log"
	"github.com/samalba/dockerclient"
	"fmt"
	"os"
	"github.com/xtraclabs/roll/rollsvcs"
	"os/signal"
)

func main() {
	//Grab the environment
	dockerHost, dockerCertPath := dockerutil.ReadDockerEnv()

	// Init the client
	log.Println("Create docker client")
	docker, _ := dockerclient.NewDockerClient(dockerHost, dockerutil.BuildDockerTLSConfig(dockerCertPath))

	containerName, token := runVault(docker)

	fmt.Printf("export VAULT_TOKEN=%s\n", token)
	fmt.Println("export VAULT_ADDR=http://localhost:8200")

	//Set up interrupt signal handler
	signalChan := make(chan os.Signal, 1)
	shutdownDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)

	//Run the server
	go func() {
		os.Setenv("VAULT_TOKEN", token)
		os.Setenv("VAULT_ADDR", "http://localhost:8200")

		coreConfig := rollsvcs.DefaultConfig()
		rollsvcs.RunRoll(3000,coreConfig)
	}()

	//Handler shutdown
	go func() {
		for _ = range signalChan {
			fmt.Println("\nReceived an interrupt, stopping roll...")
			stopVaultOnShutdown(containerName, docker)
			shutdownDone <- true
		}
	}()

	//Block until shutdown
	<-shutdownDone

}