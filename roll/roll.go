package roll

import (
	"errors"
)

//Core encapsulates the infrastructure dependencies associated with the application
type Core struct {
	developerRepo   DeveloperRepo
	applicationRepo ApplicationRepo
	SecretsRepo SecretsRepo
}

//CoreConfig is a structure used to inject infrastructure dependency implementations into
//the core struct
type CoreConfig struct {
	DeveloperRepo   DeveloperRepo
	ApplicationRepo ApplicationRepo
	SecretsRepo     SecretsRepo
}

//NewCore creates a new Core instance injecting dependencies from the CoreConfig argument
func NewCore(config *CoreConfig) *Core {

	if config.DeveloperRepo == nil {
		panic(errors.New("core config must specify a repo for developer persistance"))
	}

	if config.ApplicationRepo == nil {
		panic(errors.New("core config must specify a repo for application persistance"))
	}

	if config.SecretsRepo == nil {
		panic(errors.New("core config must specify a repo of secrets persistance"))
	}

	return &Core{
		developerRepo:   config.DeveloperRepo,
		applicationRepo: config.ApplicationRepo,
		SecretsRepo:     config.SecretsRepo,
	}
}

//StoreDeveloper stores a developer using the embedded Developer repository
func (core *Core) StoreDeveloper(dev *Developer) {
	core.developerRepo.StoreDeveloper(dev)
}

//RetrieveDeveloper retrieves a developer using the embedded Developer repository
func (core *Core) RetrieveDeveloper(email string) (*Developer, error) {
	return core.developerRepo.RetrieveDeveloper(email)
}

//StoreApplication stores an application using the embedded Application repository
func (core *Core) StoreApplication(app *Application) error {
	return core.applicationRepo.StoreApplication(app)
}

//RetrieveApplication retrieves an application using the embedded Application repository
func (core *Core) RetrieveApplication(apikey string) (*Application, error) {
	return core.applicationRepo.RetrieveApplication(apikey)
}

//StoreKeysForApp stores the private and public keys associated with an application
//using the embedded secrets store
func (core *Core) StoreKeysForApp(apikey, privateKey, publicKey string) error {
	return core.SecretsRepo.StoreKeysForApp(apikey, privateKey, publicKey)
}

//RetrievePrivateKeyForApp retrieves the private and public keys associated with an application
//using the embedded secrets store
func (core *Core) RetrievePrivateKeyForApp(apikey string) (string, error) {
	return core.SecretsRepo.RetrievePrivateKeyForApp(apikey)
}

//RetrievePublicKeyForApp retrieves the private and public keys associated with an application
//using the embedded secrets store
func (core *Core) RetrievePublicKeyForApp(apikey string) (string, error) {
	return core.SecretsRepo.RetrievePublicKeyForApp(apikey)
}
