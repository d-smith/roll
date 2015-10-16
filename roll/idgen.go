package roll

import (
	"github.com/nu7hatch/gouuid"
)

type IdGenerator interface {
	GenerateID() (string, error)
}

type UUIDIdGenerator struct{}

func (uig UUIDIdGenerator) GenerateID() (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return u.String(), nil
}
