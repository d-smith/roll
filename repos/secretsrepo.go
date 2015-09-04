package repos

import (
	"errors"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"log"
	"os"
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

//VaultSecretsRepo provides a Vault implementation of a repository for secrets. See
//https://vaultproject.io/ for more details on Vault
type VaultSecretsRepo struct {
	vaultClient *vault.Client
}

//NewVaultSecretsRepo returns a new instance of VaultSecretsRepo
func NewVaultSecretsRepo() *VaultSecretsRepo {
	vc, err := vault.NewClient(vault.DefaultConfig())
	if err != nil {
		panic(err)
	}

	return &VaultSecretsRepo{
		vaultClient: vc,
	}
}

func pathForKey(apikey string) string {
	return "secret/" + apikey
}

//StoreKeysForApp stores the private and public keys associated with an app in Vault
func (v *VaultSecretsRepo) StoreKeysForApp(apikey string, privateKey string, publicKey string) error {
	logical := v.vaultClient.Logical()
	data := make(map[string]interface{})
	data["privateKey"] = privateKey
	data["publicKey"] = publicKey
	path := pathForKey(apikey)
	s, err := logical.Write(path, data)
	if s == nil {
		log.Println("Keys for "+apikey+" written to ", path)
	}
	log.Println(fmt.Sprintf("%v", s))
	return err
}

func (v *VaultSecretsRepo) retrieveKeyFromVault(apikey string, whichKey string) (string, error) {
	logical := v.vaultClient.Logical()
	path := pathForKey(apikey)
	log.Println("Load secret from path ", path)
	secret, err := logical.Read(path)
	if err != nil {
		return "", err
	}

	if secret == nil {
		log.Println("return error - nil secret")
		return "", errors.New("No keys stored for apikey " + apikey)
	}

	var key interface{}

	switch whichKey {
	case "private":
		key = secret.Data["privateKey"]
		default:
		key = secret.Data["publicKey"]
	}

	return key.(string),nil
}

//RetrievePrivateKeyForApp retrieves the private key associated with an application  from the Vault
func (v *VaultSecretsRepo) RetrievePrivateKeyForApp(apikey string) (string, error) {
	return v.retrieveKeyFromVault(apikey, "private")
}

//RetrievePublicKeyForApp retrieves the public key associated with an application from the vault
func (v *VaultSecretsRepo) RetrievePublicKeyForApp(apikey string) (string, error) {
	return v.retrieveKeyFromVault(apikey, "public")
}
