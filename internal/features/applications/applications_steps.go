package applications

import (
	. "github.com/lsegal/gucumber"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/roll/internal/testutils"
	rollhttp "github.com/xtraclabs/roll/http"
	"net/http"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

func init() {
	var dev roll.Developer
	var app roll.Application
	var retrievedApp roll.Application
	var clientId string

	Given(`^a developer registered with the portal$`, func() {
		dev = testutils.CreateNewTestDev()
		resp := rollhttp.TestHTTPPut(T, "http://localhost:3000/v1/developers/" + dev.Email, dev)
		assert.Equal(T, http.StatusNoContent, resp.StatusCode)
	})

	And(`^they have a new application they wish to register$`, func() {
		app = roll.Application{
			ApplicationName: "int test app name",
			DeveloperEmail: dev.Email,
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

	Given(`^a registed application$`, func() {
		resp := rollhttp.TestHTTPGet(T, "http://localhost:3000/v1/applications/" + clientId, nil)
		assert.Equal(T, http.StatusOK, resp.StatusCode)

		defer resp.Body.Close()
		dec := json.NewDecoder(resp.Body)
		err := dec.Decode(&retrievedApp)
		assert.Nil(T, err)


	})

	Then(`^the details assocaited with the application can be retrieved$`, func() {
		assert.Equal(T, app.ApplicationName, retrievedApp.ApplicationName)
		assert.Equal(T, app.DeveloperEmail, retrievedApp.DeveloperEmail)
		assert.Equal(T, app.RedirectURI, retrievedApp.RedirectURI)
		assert.Equal(T, app.LoginProvider, retrievedApp.LoginProvider)
		assert.Equal(T, clientId, retrievedApp.ClientID)
		assert.True(T, len(retrievedApp.ClientSecret) > 0)
		assert.Equal(T, retrievedApp.JWTFlowPublicKey,"")
	})

}
