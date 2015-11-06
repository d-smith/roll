package roll

import (
	"errors"
)

//Core encapsulates the infrastructure dependencies associated with the application
type Core struct {
	developerRepo   DeveloperRepo
	ApplicationRepo ApplicationRepo
	AdminRepo       AdminRepo
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
	AdminRepo       AdminRepo
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
		AdminRepo:       config.AdminRepo,
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
func (core *Core) RetrieveDeveloper(email, subjectID string, adminScope bool) (*Developer, error) {
	return core.developerRepo.RetrieveDeveloper(email, subjectID, adminScope)
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

//ListDevelopers returns a list of developers registered with roll
func (core *Core) ListDevelopers(subjectID string, adminScope bool) ([]Developer, error) {
	return core.developerRepo.ListDevelopers(subjectID, adminScope)
}

//ListApplications returns a list of applications registered with roll
func (core *Core) ListApplications() ([]Application, error) {
	return core.ApplicationRepo.ListApplications()
}

//IsAdmin is a predicate used to determine if the given subject is an admin
func (core *Core) IsAdmin(subject string) (bool, error) {
	return core.AdminRepo.IsAdmin(subject)
}

//GenerateID generates and id
func (core *Core) GenerateID() (string, error) {
	return core.IdGenerator.GenerateID()
}
