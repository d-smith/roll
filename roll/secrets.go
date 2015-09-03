package roll

type SecretsRepo interface {
	StoreKeysForApp(appid string, privateKey string, publicKey string) error
}
