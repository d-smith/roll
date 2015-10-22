//This command starts vault in a docker container and spits out the VAULT_ADDR and VAULT_TOKEN to use when running
//roll integration tests
package main

import (
	"github.com/xtraclabs/roll/internal/dockerutil"
	vault "github.com/hashicorp/vault/api"
	"log"
	"net/http"
	"github.com/samalba/dockerclient"
	"time"
	"fmt"
)


//Here's the assumed docker build commands
// Vault:
//	docker build -t "vault-roll" .

const (
	VaultTestContainer      = "vault-roll"
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

func runVault(docker *dockerclient.DockerClient) string {
	var bootedContainer bool

	//Is vault running?
	log.Println("Is vault running?")
	info := dockerutil.GetAcceptanceTestContainerInfo(docker, "atest-vault")

	if info != nil {
		log.Println("Vault container found - state is: ", info.State.StateString())
		log.Fatal("You must kill and remove the container manually - can't get the root token from an existing container in this test")
		return "" //not reached
	}

	log.Println("Vault is not running - create container context")
	bootedContainer = true
	vaultContainerCtx := createVaultTestContainerContext()

	//Create and start the container.
	log.Println("Create and start the container")
	dockerutil.CreateAndStartContainer(docker, []string{"IPC_LOCK"}, nil, vaultContainerCtx)

	if bootedContainer {
		//Give the container a little time to boot
		time.Sleep(1 * time.Second)
	}

	info = dockerutil.GetAcceptanceTestContainerInfo(docker, "atest-vault")

	//Now initialize and unseal the damn vault
	rootToken := initializeVault()

	return rootToken
}

func main() {
	//Grab the environment
	dockerHost, dockerCertPath := dockerutil.ReadDockerEnv()

	// Init the client
	log.Println("Create docker client")
	docker, _ := dockerclient.NewDockerClient(dockerHost, dockerutil.BuildDockerTLSConfig(dockerCertPath))

	token := runVault(docker)

	fmt.Printf("export VAULT_TOKEN=%s\n", token)
	fmt.Println(dockerHost)
}