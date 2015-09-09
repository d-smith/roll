package roll

//SecretsRepo defines a repository abstraction for reading and writing secrets from a secure
//secret store
type SecretsRepo interface {
	StoreKeysForApp(appkey string, privateKey string, publicKey string) error
	RetrievePrivateKeyForApp(appkey string) (string, error)
	RetrievePublicKeyForApp(appkey string) (string, error)
}
