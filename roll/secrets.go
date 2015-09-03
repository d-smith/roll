package roll

type SecretsRepo interface {
	StoreKeysForApp(appkey string, privateKey string, publicKey string) error
	RetrievePrivateKeyForApp(appkey string) (string, error)
}
