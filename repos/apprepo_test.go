// +build integration

package repos

import (
	"github.com/stretchr/testify/assert"
	"github.com/xtraclabs/roll/roll"
	"strconv"
	"testing"
	"time"
)

func appPresentInList(apps []roll.Application, clientID string) bool {
	for _, app := range apps {
		if app.ClientID == clientID {
			return true
		}
	}

	return false
}

func TestListApps(t *testing.T) {

	appRepo := NewDynamoAppRepo()
	devRepo := NewDynamoDevRepo()

	testDev := testCreateDev()
	err := devRepo.StoreDeveloper(testDev)
	if !assert.Nil(t, err) {
		println(err.Error())
		t.Fail()
		return
	}

	clientID := strconv.Itoa(int(time.Now().Unix()))

	app := roll.Application{
		ClientID:        clientID,
		ClientSecret:    "xxx",
		DeveloperEmail:  testDev.Email,
		DeveloperID:     testDev.ID,
		ApplicationName: "test app",
		RedirectURI:     "http://foo.com/dev/null",
		LoginProvider:   "xtrac://loginhost:9000",
	}

	err = appRepo.CreateApplication(&app)
	if !assert.Nil(t, err) {
		t.Fail()
		return
	}

	apps, err := appRepo.ListApplications(testDev.ID, false)
	assert.Nil(t, err)
	assert.True(t, appPresentInList(apps, clientID))

	apps, err = appRepo.ListApplications("someotherid", true)
	assert.Nil(t, err)
	assert.True(t, appPresentInList(apps, clientID))

	apps, err = appRepo.ListApplications("someotherid", false)
	assert.Nil(t, err)
	assert.True(t, len(apps) == 0)

}

func TestUpdateApplication(t *testing.T) {
	appRepo := NewDynamoAppRepo()
	devRepo := NewDynamoDevRepo()

	testDev := testCreateDev()
	err := devRepo.StoreDeveloper(testDev)
	if !assert.Nil(t, err) {
		println(err.Error())
		t.Fail()
		return
	}

	clientID := "a-" + strconv.Itoa(int(time.Now().Unix()))
	appName := "app" + strconv.Itoa(int(time.Now().Unix()))
	updatedAppName := "update-app" + strconv.Itoa(int(time.Now().Unix()))

	app := roll.Application{
		ClientID:        clientID,
		ClientSecret:    "xxx",
		DeveloperEmail:  testDev.Email,
		DeveloperID:     testDev.ID,
		ApplicationName: appName,
		RedirectURI:     "http://foo.com/dev/null",
		LoginProvider:   "xtrac://loginhost:9000",
	}

	err = appRepo.CreateApplication(&app)
	if !assert.Nil(t, err) {
		t.Fail()
		return
	}

	retrieved, err := appRepo.RetrieveApplication(clientID)
	if !assert.Nil(t, err) {
		t.Fail()
		return
	}

	assert.Equal(t, app.ClientID, retrieved.ClientID)
	assert.Equal(t, app.ClientSecret, retrieved.ClientSecret)
	assert.Equal(t, app.DeveloperEmail, retrieved.DeveloperEmail)
	assert.Equal(t, app.DeveloperID, retrieved.DeveloperID)
	assert.Equal(t, app.ApplicationName, retrieved.ApplicationName)
	assert.Equal(t, app.RedirectURI, retrieved.RedirectURI)
	assert.Equal(t, app.LoginProvider, retrieved.LoginProvider)

	t.Log("Update when not the owner generates an error")
	retrieved.ApplicationName = updatedAppName
	err = appRepo.UpdateApplication(retrieved, "not the owner")
	if !assert.NotNil(t, err) {
		t.Fail()
		return
	}

	_, ok := err.(roll.NonOwnerUpdateError)
	if !assert.True(t, ok) {
		t.Fail()
		return
	}

	t.Log("Update application as owner succeeds")
	err = appRepo.UpdateApplication(retrieved, testDev.ID)
	if !assert.Nil(t, err) {
		t.Fail()
		return
	}

	updated, err := appRepo.RetrieveApplication(clientID)
	if !assert.Nil(t, err) {
		t.Fail()
		return
	}

	assert.Equal(t, retrieved.ClientID, updated.ClientID)
	assert.Equal(t, retrieved.ClientSecret, updated.ClientSecret)
	assert.Equal(t, retrieved.DeveloperEmail, updated.DeveloperEmail)
	assert.Equal(t, retrieved.DeveloperID, updated.DeveloperID)
	assert.Equal(t, retrieved.ApplicationName, updated.ApplicationName)
	assert.Equal(t, retrieved.RedirectURI, updated.RedirectURI)
	assert.Equal(t, retrieved.LoginProvider, updated.LoginProvider)

}
