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

	var devsCountAdmin int

	var email = "foo@foo.com"
	dev := roll.Developer{
		Email:     email,
		FirstName: "Foo",
		LastName:  "Barr",
		ID:        "foo",
	}

	devRepo := NewMBDDevRepo()

	//Grab the baseline count of all users from the perspective of an admin user. We expect
	//the admin count + 1 at the end of the test
	devs, err := devRepo.ListDevelopers("no one", true)
	assert.Nil(t, err)
	devsCountAdmin = len(devs)

	devs, err = devRepo.ListDevelopers("no one", false)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(devs))

	err = devRepo.StoreDeveloper(&dev)
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

	devs, err = devRepo.ListDevelopers("no one", true)
	assert.Nil(t, err)
	assert.Equal(t, devsCountAdmin+1, len(devs))

	devs, err = devRepo.ListDevelopers("no one", false)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(devs))

	devs, err = devRepo.ListDevelopers(dev.ID, false)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(devs))

}
