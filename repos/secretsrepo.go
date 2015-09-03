package repos
import "errors"


type VaultSecretsRepo struct {}


func NewVaultSecretsRepo() *VaultSecretsRepo {
	return &VaultSecretsRepo{}
}

func (v *VaultSecretsRepo) StoreKeysForApp(appid string, privateKey string, publicKey string) error {
	return errors.New("Not Implemented")
}
