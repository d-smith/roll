// +build integration

package mdb

import (
	"github.com/stretchr/testify/assert"
	"github.com/xtraclabs/roll/roll"
	"testing"
)

func TestRetrieveNonexistentApp(t *testing.T) {
	var app *roll.Application
	appRepo := NewMBDAppRepo()
	app, err := appRepo.RetrieveAppByNameAndDevEmail("xx", "yy")
	assert.NotNil(t, err)
	assert.Nil(t, app)
}

func TestAddAndRetrieveApp(t *testing.T) {
	app := new(roll.Application)
	app.ApplicationName = "an app"
	app.ClientID = "123"
	app.ClientSecret = "hush"
	app.DeveloperEmail = "foo@foo.bar"
	app.DeveloperID = "foo"
	app.LoginProvider = "auth0"
	app.RedirectURI = "neither here nor there"

	appRepo := NewMBDAppRepo()
	err := appRepo.CreateApplication(app)
	if assert.Nil(t, err) {
		defer appRepo.delete(app)
	}

	retapp, err := appRepo.RetrieveAppByNameAndDevEmail("an app", "foo@foo.bar")
	assert.Nil(t, err)
	if assert.NotNil(t, app) {
		assert.Equal(t, app.ApplicationName, retapp.ApplicationName)
		assert.Equal(t, app.ClientID, retapp.ClientID)
		assert.Equal(t, app.ClientSecret, retapp.ClientSecret)
		assert.Equal(t, app.DeveloperEmail, retapp.DeveloperEmail)
		assert.Equal(t, app.DeveloperID, retapp.DeveloperID)
		assert.Equal(t, app.LoginProvider, retapp.LoginProvider)
		assert.Equal(t, app.RedirectURI, retapp.RedirectURI)
	}

	retapp, err = appRepo.RetrieveApplication(app.ClientID, app.DeveloperID, false)
	assert.Nil(t, err)
	if assert.NotNil(t, app) {
		assert.Equal(t, app.ApplicationName, retapp.ApplicationName)
		assert.Equal(t, app.ClientID, retapp.ClientID)
		assert.Equal(t, app.ClientSecret, retapp.ClientSecret)
		assert.Equal(t, app.DeveloperEmail, retapp.DeveloperEmail)
		assert.Equal(t, app.DeveloperID, retapp.DeveloperID)
		assert.Equal(t, app.LoginProvider, retapp.LoginProvider)
		assert.Equal(t, app.RedirectURI, retapp.RedirectURI)
	}

	retapp, err = appRepo.RetrieveApplication(app.ClientID, "huh", true)
	assert.Nil(t, err)
	if assert.NotNil(t, app) {
		assert.Equal(t, app.ApplicationName, retapp.ApplicationName)
		assert.Equal(t, app.ClientID, retapp.ClientID)
		assert.Equal(t, app.ClientSecret, retapp.ClientSecret)
		assert.Equal(t, app.DeveloperEmail, retapp.DeveloperEmail)
		assert.Equal(t, app.DeveloperID, retapp.DeveloperID)
		assert.Equal(t, app.LoginProvider, retapp.LoginProvider)
		assert.Equal(t, app.RedirectURI, retapp.RedirectURI)
	}

	retapp, err = appRepo.SystemRetrieveApplication(app.ClientID)
	assert.Nil(t, err)
	assert.Equal(t, app.ClientID, retapp.ClientID)

	retapp, err = appRepo.RetrieveApplication(app.ClientID, "huh", false)
	assert.NotNil(t, err)
	assert.Nil(t, retapp)
}

func TestSecretGeneratedWhenNeede(t *testing.T) {
	app := new(roll.Application)
	app.ApplicationName = "an app"
	app.ClientID = "123"
	app.DeveloperEmail = "foo@foo.bar"
	app.DeveloperID = "foo"
	app.LoginProvider = "auth0"
	app.RedirectURI = "neither here nor there"

	appRepo := NewMBDAppRepo()
	err := appRepo.CreateApplication(app)
	if assert.Nil(t, err) {
		defer appRepo.delete(app)
	}

	retapp, err := appRepo.RetrieveAppByNameAndDevEmail("an app", "foo@foo.bar")
	assert.Nil(t, err)
	assert.NotEqual(t, "", retapp.ClientSecret)
}

func TestDuplicateAppCreateGeneratesError(t *testing.T) {
	app := new(roll.Application)
	app.ApplicationName = "an app"
	app.ClientID = "123"
	app.DeveloperEmail = "foo@foo.bar"
	app.DeveloperID = "foo"
	app.LoginProvider = "auth0"
	app.RedirectURI = "neither here nor there"

	appRepo := NewMBDAppRepo()
	err := appRepo.CreateApplication(app)
	if assert.Nil(t, err) {
		defer appRepo.delete(app)
	}

	err = appRepo.CreateApplication(app)
	assert.NotNil(t, err)
}

func TestUpdateApp(t *testing.T) {

	appRepo := NewMBDAppRepo()

	//Count the apps prior to creating one
	apps, err := appRepo.ListApplications("foo", true)
	assert.Nil(t, err)
	adminCount := len(apps)

	//No apps see with a user id of not foo and not an admin
	apps, err = appRepo.ListApplications("not foo", false)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(apps))

	//Create an app
	app := new(roll.Application)
	app.ApplicationName = "an app"
	app.ClientID = "123"
	app.DeveloperEmail = "foo@foo.bar"
	app.DeveloperID = "foo"
	app.LoginProvider = "auth0"
	app.RedirectURI = "neither here nor there"

	err = appRepo.CreateApplication(app)
	if assert.Nil(t, err) {
		defer appRepo.delete(app)
	}

	err = appRepo.UpdateApplication(app, "no way jose")
	assert.NotNil(t, err)

	err = appRepo.UpdateApplication(app, app.DeveloperID)
	assert.Nil(t, err)

	app.JWTFlowAudience = "aud"
	app.JWTFlowIssuer = "iss"
	app.JWTFlowPublicKey = "key to the city"
	appRepo.UpdateApplication(app, app.DeveloperID)

	retapp, err := appRepo.SystemRetrieveApplicationByJWTFlowAudience("aud")
	assert.Nil(t, err)
	if assert.NotNil(t, app) {
		assert.Equal(t, app.ApplicationName, retapp.ApplicationName)
		assert.Equal(t, app.ClientID, retapp.ClientID)
		assert.Equal(t, app.ClientSecret, retapp.ClientSecret)
		assert.Equal(t, app.DeveloperEmail, retapp.DeveloperEmail)
		assert.Equal(t, app.DeveloperID, retapp.DeveloperID)
		assert.Equal(t, app.LoginProvider, retapp.LoginProvider)
		assert.Equal(t, app.RedirectURI, retapp.RedirectURI)
		assert.Equal(t, app.JWTFlowAudience, retapp.JWTFlowAudience)
		assert.Equal(t, app.JWTFlowIssuer, retapp.JWTFlowIssuer)
		assert.Equal(t, app.JWTFlowPublicKey, retapp.JWTFlowPublicKey)
	}

	//Admin user should see an additional app in the list
	apps, err = appRepo.ListApplications("foo", true)
	assert.Nil(t, err)
	assert.Equal(t, adminCount+1, len(apps))

	//User adding the app should see a list with 1 entry
	apps, err = appRepo.ListApplications("foo", false)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(apps))
}

func TestUpdateNoSuchApp(t *testing.T) {

	appRepo := NewMBDAppRepo()

	//Specify an app
	app := new(roll.Application)
	app.ApplicationName = "an app"
	app.ClientID = "123"
	app.DeveloperEmail = "foo@foo.bar"
	app.DeveloperID = "foo"
	app.LoginProvider = "auth0"
	app.RedirectURI = "neither here nor there"

	err := appRepo.UpdateApplication(app, app.DeveloperID)
	assert.NotNil(t, err)
}
