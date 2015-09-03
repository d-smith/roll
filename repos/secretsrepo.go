package repos
import (
	"errors"
	"os"
	vault "github.com/hashicorp/vault/api"
	"fmt"
	"log"
)

var vaultEndpoint string
var vaultToken string

func init() {
	vaultEndpoint = os.Getenv("VAULT_ADDR")
	if vaultEndpoint == "" {
		panic(errors.New("Missing Configuration: VAULT_ADDR env variable not specified"))
	}

	vaultToken = os.Getenv("VAULT_TOKEN")
	if vaultToken == "" {
		panic(errors.New("Missing configuration: VAULT_TOKEN env variable not specified"))
	}
}

type VaultSecretsRepo struct {
	vaultClient *vault.Client
}


func NewVaultSecretsRepo() *VaultSecretsRepo {
	vc, err := vault.NewClient(vault.DefaultConfig())
	if err != nil {
		panic(err)
	}

	return &VaultSecretsRepo{
		vaultClient: vc,
	}
}

func (v *VaultSecretsRepo) StoreKeysForApp(apikey string, privateKey string, publicKey string) error {
	logical := v.vaultClient.Logical()
	data := make(map[string]interface{})
	data["privateKey"] = privateKey
	data["publicKey"] = publicKey
	path := "secret/" + apikey
	s, err := logical.Write(path, data)
	if s == nil {
		log.Println("Keys for " + apikey + " written to ", path)
	}
	log.Println(fmt.Sprintf("%v", s))
	return err
}
