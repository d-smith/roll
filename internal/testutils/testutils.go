package testutils

import (
	"fmt"
	"time"
	"github.com/xtraclabs/roll/roll"
)

func NewDevTestEmail() string {
	return fmt.Sprintf("%d@%d.net", time.Now().Unix(), time.Now().Unix())
}

func CreateNewTestDev() roll.Developer {
	return roll.Developer{
		FirstName: "test",
		LastName: "test",
		Email: NewDevTestEmail(),
	}
}
