package roll

import (
	"errors"
)

type Core struct {
	developerRepo DeveloperRepo
	applicationRepo ApplicationRepo
}

type CoreConfig struct {
	DeveloperRepo   DeveloperRepo
	ApplicationRepo ApplicationRepo
}

func NewCore(config *CoreConfig) *Core {

	if config.DeveloperRepo == nil {
		panic(errors.New("core config must specify a repo for developer persistance"))
	}

	if config.ApplicationRepo == nil {
		panic(errors.New("core config must specify a repo for application persistance"))
	}


	return &Core{
		developerRepo: config.DeveloperRepo,
		applicationRepo: config.ApplicationRepo,
	}
}

func (core *Core) StoreDeveloper(dev *Developer) {
	core.developerRepo.StoreDeveloper(dev)
}

func (core *Core) RetrieveDeveloper(email string) (*Developer, error) {
	return core.developerRepo.RetrieveDeveloper(email)
}

func (core *Core) StoreApplication(app *Application) error {
	return core.applicationRepo.StoreApplication(app)
}

func (core *Core) RetrieveApplication(apikey string) (*Application, error) {
	return core.applicationRepo.RetrieveApplication(apikey)
}
