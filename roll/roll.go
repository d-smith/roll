package roll

import (
	"errors"
)

//Core encapsulates the infrastructure dependencies associated with the application
type Core struct {
	developerRepo   DeveloperRepo
	ApplicationRepo ApplicationRepo
	SecretsRepo     SecretsRepo
	IdGenerator     IdGenerator
	secure          bool
	rollClientId    string
}

//CoreConfig is a structure used to inject infrastructure dependency implementations into
//the core struct
type CoreConfig struct {
	DeveloperRepo   DeveloperRepo
	ApplicationRepo ApplicationRepo
	SecretsRepo     SecretsRepo
	IdGenerator     IdGenerator
	Secure          bool
	RollClientID    string
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
		panic(errors.New("core config must specify a repo for secrets persistance"))
	}

	if config.IdGenerator == nil {
		panic(errors.New("core config must specify an id generator"))
	}

	return &Core{
		developerRepo:   config.DeveloperRepo,
		ApplicationRepo: config.ApplicationRepo,
		SecretsRepo:     config.SecretsRepo,
		IdGenerator:     config.IdGenerator,
		secure:          config.Secure,
		rollClientId:    config.RollClientID,
	}
}

//Secure returns true if roll is running in secure mode, false otherwise
func (core *Core) Secure() bool {
	return core.secure
}

//StoreDeveloper stores a developer using the embedded Developer repository
func (core *Core) StoreDeveloper(dev *Developer) error {
	return core.developerRepo.StoreDeveloper(dev)
}

//RetrieveDeveloper retrieves a developer using the embedded Developer repository
func (core *Core) RetrieveDeveloper(email string) (*Developer, error) {
	return core.developerRepo.RetrieveDeveloper(email)
}

//StoreApplication stores an application using the embedded Application repository
func (core *Core) CreateApplication(app *Application) error {
	return core.ApplicationRepo.CreateApplication(app)
}

//StoreApplication stores an application using the embedded Application repository
func (core *Core) UpdateApplication(app *Application) error {
	return core.ApplicationRepo.UpdateApplication(app)
}

//RetrieveApplication retrieves an application using the embedded Application repository
func (core *Core) RetrieveApplication(clientID string) (*Application, error) {
	return core.ApplicationRepo.RetrieveApplication(clientID)
}

//StoreKeysForApp stores the private and public keys associated with an application
//using the embedded secrets store
func (core *Core) StoreKeysForApp(clientID, privateKey, publicKey string) error {
	return core.SecretsRepo.StoreKeysForApp(clientID, privateKey, publicKey)
}

//RetrievePrivateKeyForApp retrieves the private and public keys associated with an application
//using the embedded secrets store
func (core *Core) RetrievePrivateKeyForApp(clientID string) (string, error) {
	return core.SecretsRepo.RetrievePrivateKeyForApp(clientID)
}

//RetrievePublicKeyForApp retrieves the private and public keys associated with an application
//using the embedded secrets store
func (core *Core) RetrievePublicKeyForApp(clientID string) (string, error) {
	return core.SecretsRepo.RetrievePublicKeyForApp(clientID)
}

func (core *Core) ListDevelopers() ([]Developer, error) {
	return core.developerRepo.ListDevelopers()
}

func (core *Core) ListApplications() ([]Application, error) {
	return core.ApplicationRepo.ListApplications()
}

func (core *Core) GenerateID() (string, error) {
	return core.IdGenerator.GenerateID()
}
