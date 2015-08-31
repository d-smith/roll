package roll

import (
	"errors"
)

type Core struct {
	developerRepo DeveloperRepo
}

type CoreConfig struct {
	DeveloperRepo DeveloperRepo
}

func NewCore(config *CoreConfig) *Core {
	if config.DeveloperRepo == nil {
		panic(errors.New("core config must specify a repo for developer persistance"))
	}
	return &Core{
		developerRepo: config.DeveloperRepo,
	}
}

func (core *Core) StoreDeveloper(dev *Developer) {
	core.developerRepo.StoreDeveloper(dev)
}

func (core *Core) RetrieveDeveloper(email string) (*Developer, error) {
	return core.developerRepo.RetrieveDeveloper(email)
}
