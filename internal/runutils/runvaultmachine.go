//This command starts vault in a docker container and spits out the VAULT_ADDR and VAULT_TOKEN to use when running
//roll integration tests
package runutils

import (
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"github.com/samalba/dockerclient"
	"github.com/xtraclabs/roll/internal/dockerutil"
	"github.com/xtraclabs/roll/rollsvcs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

//Here's the assumed docker build commands
// Vault:
//	docker build -t "vault-roll" .

const (
	VaultTestContainer = "vault-roll"
)

func createVaultTestContainerContext() *dockerutil.ContainerContext {
	containerCtx := dockerutil.ContainerContext{
		ImageName: VaultTestContainer,
	}

	containerCtx.Labels = make(map[string]string)
	containerCtx.Labels["xt-container-type"] = "atest-vault"

	containerCtx.PortContext = make(map[string]string)
	containerCtx.PortContext["8200/tcp"] = "8200"

	return &containerCtx
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func createVaultClient() *vault.Client {
	log.Println("Create vault client")
	config := &vault.Config{
		Address:    "http://localhost:8200",
		HttpClient: http.DefaultClient,
	}
	vc, err := vault.NewClient(config)
	fatal(err)

	return vc
}

func unsealVault(vc *vault.Client, initResponse *vault.InitResponse) string {
	log.Println("Unseal vault")
	_, err := vc.Sys().Unseal(initResponse.Keys[0])
	fatal(err)
	return initResponse.RootToken
}

func initializeNewVault(vc *vault.Client) *vault.InitResponse {
	log.Println("Initialize fresh vault")
	vaultInit := &vault.InitRequest{
		SecretShares:    1,
		SecretThreshold: 1,
	}
	initResponse, err := vc.Sys().Init(vaultInit)
	fatal(err)

	return initResponse

}

func initializeVault() string {
	vc := createVaultClient()
	initResponse := initializeNewVault(vc)
	return unsealVault(vc, initResponse)
}

func runVault(docker *dockerclient.DockerClient) (string, string) {
	var bootedContainer bool

	//Is vault running?
	log.Println("Is vault running?")
	info := dockerutil.GetAcceptanceTestContainerInfo(docker, "atest-vault")

	if info != nil {
		log.Println("Vault container found - state is: ", info.State.StateString())
		log.Fatal("You must kill and remove the container manually - can't get the root token from an existing container in this test")
		return "", "" //not reached
	}

	log.Println("Vault is not running - create container context")
	bootedContainer = true
	vaultContainerCtx := createVaultTestContainerContext()

	//Create and start the container.
	log.Println("Create and start the container")
	containerId := dockerutil.CreateAndStartContainer(docker, []string{"IPC_LOCK"}, nil, vaultContainerCtx)

	if bootedContainer {
		//Give the container a little time to boot
		time.Sleep(1 * time.Second)
	}

	//Now initialize and unseal the damn vault
	rootToken := initializeVault()

	return containerId, rootToken
}

func stopVaultOnShutdown(containerId string, docker *dockerclient.DockerClient) {
	log.Println("... stopping container", containerId, "...")
	docker.StopContainer(containerId, 5)
	docker.RemoveContainer(containerId, true, false)
}

func RunVaultAndRoll() chan bool {
	//Grab the environment
	dockerHost, dockerCertPath := dockerutil.ReadDockerEnv()

	// Init the client
	log.Println("Create docker client")
	docker, _ := dockerclient.NewDockerClient(dockerHost, dockerutil.BuildDockerTLSConfig(dockerCertPath))

	containerName, token := runVault(docker)

	log.Printf("export VAULT_TOKEN=%s\n", token)
	log.Println("export VAULT_ADDR=http://localhost:8200")

	//Set up interrupt signal handler
	signalChan := make(chan os.Signal, 1)
	shutdownDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)

	//Run the server
	go func() {
		os.Setenv("VAULT_TOKEN", token)
		os.Setenv("VAULT_ADDR", "http://localhost:8200")

		coreConfig := rollsvcs.DefaultUnsecureConfig()
		rollsvcs.RunRoll(3000, coreConfig)
	}()

	//Handler shutdown
	go func() {
		for _ = range signalChan {
			fmt.Println("\nReceived an interrupt, stopping roll...")
			stopVaultOnShutdown(containerName, docker)
			shutdownDone <- true
		}
	}()

	return shutdownDone
}
