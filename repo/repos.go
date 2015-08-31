package repo

import (
	"github.com/xtraclabs/roll/roll"
)

type DeveloperRepo interface {
	RetrieveDeveloper() (*roll.Developer, error)
	StoreDeveloper(*roll.Developer, error)
}
