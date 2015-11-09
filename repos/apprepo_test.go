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
