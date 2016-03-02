package repos

import (
	"errors"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	log "github.com/Sirupsen/logrus"
	"os"
)

//VaultSecretsRepo provides a Vault implementation of a repository for secrets. See
//https://vaultproject.io/ for more details on Vault
type VaultSecretsRepo struct {
	vaultEndpoint string
	vaultToken    string
	vaultClient   *vault.Client
}

//NewVaultSecretsRepo returns a new instance of VaultSecretsRepo
func NewVaultSecretsRepo() *VaultSecretsRepo {
	vaultEndpoint := os.Getenv("VAULT_ADDR")
	if vaultEndpoint == "" {
		panic(errors.New("Missing Configuration: VAULT_ADDR env variable not specified"))
	}

	vaultToken := os.Getenv("VAULT_TOKEN")
	if vaultToken == "" {
		panic(errors.New("Missing configuration: VAULT_TOKEN env variable not specified"))
	}

	vc, err := vault.NewClient(vault.DefaultConfig())
	if err != nil {
		panic(err)
	}

	return &VaultSecretsRepo{
		vaultClient:   vc,
		vaultEndpoint: vaultEndpoint,
		vaultToken:    vaultToken,
	}
}

func pathForKey(clientID string) string {
	return "secret/" + clientID
}

//StoreKeysForApp stores the private and public keys associated with an app in Vault
func (v *VaultSecretsRepo) StoreKeysForApp(clientID string, privateKey string, publicKey string) error {
	logical := v.vaultClient.Logical()
	data := make(map[string]interface{})
	data["privateKey"] = privateKey
	data["publicKey"] = publicKey
	path := pathForKey(clientID)
	s, err := logical.Write(path, data)
	if s == nil {
		log.Info("Keys for "+clientID+" written to ", path)
	}
	log.Info(fmt.Sprintf("%v", s))
	return err
}

func (v *VaultSecretsRepo) retrieveKeyFromVault(clientID string, whichKey string) (string, error) {
	logical := v.vaultClient.Logical()
	path := pathForKey(clientID)
	log.Info("Load secret from path ", path)
	secret, err := logical.Read(path)
	if err != nil {
		return "", err
	}

	if secret == nil {
		log.Info("return error - nil secret")
		return "", errors.New("No keys stored for clientID " + clientID)
	}

	var key interface{}

	switch whichKey {
	case "private":
		key = secret.Data["privateKey"]
	default:
		key = secret.Data["publicKey"]
	}

	return key.(string), nil
}

//RetrievePrivateKeyForApp retrieves the private key associated with an application  from the Vault
func (v *VaultSecretsRepo) RetrievePrivateKeyForApp(clientID string) (string, error) {
	return v.retrieveKeyFromVault(clientID, "private")
}

//RetrievePublicKeyForApp retrieves the public key associated with an application from the vault
func (v *VaultSecretsRepo) RetrievePublicKeyForApp(clientID string) (string, error) {
	return v.retrieveKeyFromVault(clientID, "public")
}
