package applications

import (
	"encoding/json"
	. "github.com/lsegal/gucumber"
	"github.com/stretchr/testify/assert"
	rollhttp "github.com/xtraclabs/roll/http"
	"github.com/xtraclabs/roll/internal/testutils"
	"github.com/xtraclabs/roll/roll"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func init() {
	var dev roll.Developer
	var app roll.Application
	var retrievedApp roll.Application
	var clientId string
	var reRegisterStatus int
	var duplicationErrorMessage string

	Before("@apptests", func() {
		testutils.URLGuard("http://localhost:3000/v1/developers")
	})

	Given(`^a developer registered with the portal$`, func() {
		dev = testutils.CreateNewTestDev()
		resp := rollhttp.TestHTTPPut(T, "http://localhost:3000/v1/developers/"+dev.Email, dev)
		println("resp is", resp)
		assert.Equal(T, http.StatusNoContent, resp.StatusCode)
	})

	And(`^they have a new application they wish to register$`, func() {
		app = roll.Application{
			ApplicationName: "int test app name",
			DeveloperEmail:  dev.Email,
			RedirectURI:     "http://localhost:3000/ab",
			LoginProvider:   "xtrac://localhost:9000",
		}
	})

	Then(`^the application should be successfully registered$`, func() {
		resp := rollhttp.TestHTTPPost(T, "http://localhost:3000/v1/applications", app)
		assert.Equal(T, http.StatusOK, resp.StatusCode)

		var appCreatedResponse rollhttp.ApplicationCreatedResponse

		dec := json.NewDecoder(resp.Body)
		err := dec.Decode(&appCreatedResponse)
		assert.Nil(T, err)
		assert.True(T, len(appCreatedResponse.ClientID) > 0)
		clientId = appCreatedResponse.ClientID
	})

	Given(`^a registered application$`, func() {
		retrieveAppDefinition(clientId, &retrievedApp)
	})

	Then(`^the details associated with the application can be retrieved$`, func() {
		assert.Equal(T, app.ApplicationName, retrievedApp.ApplicationName)
		assert.Equal(T, app.DeveloperEmail, retrievedApp.DeveloperEmail)
		assert.Equal(T, app.RedirectURI, retrievedApp.RedirectURI)
		assert.Equal(T, app.LoginProvider, retrievedApp.LoginProvider)
		assert.Equal(T, clientId, retrievedApp.ClientID)
		assert.True(T, len(retrievedApp.ClientSecret) > 0)
		assert.Equal(T, retrievedApp.JWTFlowPublicKey, "")
	})

	Given(`^an application has already been registered$`, func() {
		assert.True(T, len(clientId) > 0)
	})

	And(`^a developer attempts to register an application with the same name$`, func() {
		resp := rollhttp.TestHTTPPost(T, "http://localhost:3000/v1/applications", app)
		reRegisterStatus = resp.StatusCode

		defer resp.Body.Close()
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.Nil(T, err)
		duplicationErrorMessage = string(bodyBytes)
	})

	Then(`^an error is returned with status code StatusConflict$`, func() {
		assert.Equal(T, http.StatusConflict, reRegisterStatus)
	})

	And(`^the error message indicates a duplicate registration was attempted$`, func() {
		assert.True(T, strings.Contains(duplicationErrorMessage, "definition exists for application"))
	})

	Given(`^a registered application to update$`, func() {
		assert.True(T, len(clientId) > 0)
	})

	And(`^there are updates to make to the application defnition$`, func() {
		app.RedirectURI = "http://localhost:3000/son_of_callback"
	})

	Then(`^the application can be updated$`, func() {
		resp := rollhttp.TestHTTPPut(T, "http://localhost:3000/v1/applications/"+clientId, app)
		assert.Equal(T, http.StatusNoContent, resp.StatusCode)
	})

	And(`^the updates are reflected when retrieving the application definition anew$`, func() {
		retrieveAppDefinition(clientId, &retrievedApp)
		assert.Equal(T, "http://localhost:3000/son_of_callback", retrievedApp.RedirectURI)
	})

}

func retrieveAppDefinition(clientID string, app interface{}) {
	resp := rollhttp.TestHTTPGet(T, "http://localhost:3000/v1/applications/"+clientID, nil)
	assert.Equal(T, http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&app)
	assert.Nil(T, err)
	log.Printf("%v\n", app)
}
