package http

import (
	"github.com/xtraclabs/roll/roll"
)

func grantAdminScope(core *roll.Core, subject string) (bool, error) {
	return core.IsAdmin(subject)
}
