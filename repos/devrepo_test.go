// +build integration

package repos

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/xtraclabs/roll/roll"
	"strconv"
	"testing"
	"time"
)

func devPresent(id string, devs []roll.Developer) bool {
	for _, dev := range devs {
		if dev.ID == id {
			return true
		}
	}

	return false
}

func TestDevRepoStuff(t *testing.T) {
	devrepo := NewDynamoDevRepo()

	t.Log("Store a developer")

	testDevID := strconv.Itoa(int(time.Now().Unix()))
	firstName := testDevID + "fn"
	lastName := testDevID + "ln"
	testEmail := testDevID + "@foo.com"

	validateRetrieveAttrs := func(retDev *roll.Developer) {
		assert.Equal(t, testDevID, retDev.ID)
		assert.Equal(t, firstName, retDev.FirstName)
		assert.Equal(t, lastName, retDev.LastName)
		assert.Equal(t, testEmail, retDev.Email)
	}

	dev := roll.Developer{
		ID:        testDevID,
		Email:     testEmail,
		FirstName: firstName,
		LastName:  lastName,
	}

	err := devrepo.StoreDeveloper(&dev)
	assert.Nil(t, err)

	t.Log("List developers")
	devs, err := devrepo.ListDevelopers("xxx", true)
	assert.Nil(t, err)
	assert.True(t, devPresent(testDevID, devs))

	t.Log("List developers with matching subject id")
	devs, err = devrepo.ListDevelopers(testDevID, false)
	if !assert.Nil(t, err) {
		println(err.Error())
	}
	fmt.Println(devs)
	assert.True(t, devPresent(testDevID, devs))

	t.Log("List developers with no subject id match")
	t.Log("List developers with matching subject id")
	devs, err = devrepo.ListDevelopers("xxx", false)
	if !assert.Nil(t, err) {
		println(err.Error())
	}
	assert.Equal(t, 0, len(devs))

	t.Log("retrieve the developer stored earlier with no subject id restriction")
	rd, err := devrepo.RetrieveDeveloper(testEmail, "xxx", true)
	if !assert.Nil(t, err) {
		println(err.Error())
	}
	if assert.NotNil(t, rd) {
		validateRetrieveAttrs(rd)
	}

	t.Log("retrieve the developer with a subject id restriction that is satisfied")
	rd, err = devrepo.RetrieveDeveloper(testEmail, testDevID, false)
	if !assert.Nil(t, err) {
		println(err.Error())
	}
	if assert.NotNil(t, rd) {
		validateRetrieveAttrs(rd)
	}

	t.Log("verify no data returned if subject id doesn't match record")
	rd, err = devrepo.RetrieveDeveloper(testEmail, "xxx", false)
	assert.Nil(t, err)
	assert.Nil(t, rd)

}
