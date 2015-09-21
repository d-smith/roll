package http

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/roll/roll/mocks"
	"net/http"
	"testing"
)

func TestStoreApp(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	app := roll.Application{
		ApplicationName: "ambivilant birds",
		DeveloperEmail:  "doug@dev.com",
		CLientID:        "1111-2222-3333333-4444444",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("StoreApplication", &app).Return(nil)

	secretsRepoMock := coreConfig.SecretsRepo.(*mocks.SecretsRepo)
	secretsRepoMock.On("StoreKeysForApp",
		mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()

	resp := testHTTPPut(t, addr+"/v1/applications/1111-2222-3333333-4444444", app)
	appRepoMock.AssertCalled(t, "StoreApplication", &app)
	secretsRepoMock.AssertExpectations(t)

	checkResponseStatus(t, resp, http.StatusNoContent)
}

func TestGetApplication(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	returnVal := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		CLientID:        "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		ClientSecret:    "not for browser clients",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	resp := testHTTPGet(t, addr+"/v1/applications/1111-2222-3333333-4444444", nil)
	appRepoMock.AssertCalled(t, "RetrieveApplication", "1111-2222-3333333-4444444")

	var actual roll.Application

	checkResponseBody(t, resp, &actual)
	assert.Equal(t, "doug@dev.com", actual.DeveloperEmail)
	assert.Equal(t, "1111-2222-3333333-4444444", actual.CLientID)
	assert.Equal(t, "fight club", actual.ApplicationName)
	assert.Equal(t, "not for browser clients", actual.ClientSecret)
	assert.Equal(t, "http://localhost:3000/ab", actual.RedirectURI)
	assert.Equal(t, "xtrac://localhost:9000", actual.LoginProvider)

}
