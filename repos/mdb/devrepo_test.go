// +build integration

package mdb

import (
	"github.com/stretchr/testify/assert"
	"github.com/xtraclabs/roll/roll"
	"testing"
)

func TestRetrieveNonExistentDev(t *testing.T) {
	devRepo := NewMBDDevRepo()
	_, err := devRepo.RetrieveDeveloper("nothing", "foo", true)
	assert.NotNil(t, err)

	_, err = devRepo.RetrieveDeveloper("nothing", "foo", false)
	assert.NotNil(t, err)
}

func TestCreateAndRetrieveDev(t *testing.T) {
	var email = "foo@foo.com"
	dev := roll.Developer{
		Email:     email,
		FirstName: "Foo",
		LastName:  "Barr",
		ID:        "foo",
	}

	devRepo := NewMBDDevRepo()
	err := devRepo.StoreDeveloper(&dev)
	assert.Nil(t, err)

	defer devRepo.deleteDeveloper(email)

	retDev, err := devRepo.RetrieveDeveloper(dev.Email, "foo", false)
	assert.Nil(t, err)
	assert.Equal(t, dev.Email, retDev.Email)
	assert.Equal(t, dev.FirstName, retDev.FirstName)
	assert.Equal(t, dev.LastName, retDev.LastName)
	assert.Equal(t, dev.ID, retDev.ID)

	retDev, err = devRepo.RetrieveDeveloper(dev.Email, "notfoo", true)
	assert.Nil(t, err)
	assert.Equal(t, dev.Email, retDev.Email)
	assert.Equal(t, dev.FirstName, retDev.FirstName)
	assert.Equal(t, dev.LastName, retDev.LastName)
	assert.Equal(t, dev.ID, retDev.ID)

}
