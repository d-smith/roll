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
	})

}
