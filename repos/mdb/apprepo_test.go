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

	apps, err := appRepo.ListApplications("foo", true)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(apps))
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

	retapp, err = appRepo.RetrieveApplication(app.ClientID, "huh", false)
	assert.NotNil(t, err)
	assert.Nil(t, retapp)
}
