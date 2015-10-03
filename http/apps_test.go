package http

import (
	"bytes"
	"encoding/json"
	"errors"
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
		ClientID:        "1111-2222-3333333-4444444",
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

func TestStoreAppSecretStoreFault(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	app := roll.Application{
		ApplicationName: "ambivilant birds",
		DeveloperEmail:  "doug@dev.com",
		ClientID:        "1111-2222-3333333-4444444",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	secretsRepoMock := coreConfig.SecretsRepo.(*mocks.SecretsRepo)
	secretsRepoMock.On("StoreKeysForApp",
		mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New("secret store error")).Once()

	resp := testHTTPPut(t, addr+"/v1/applications/1111-2222-3333333-4444444", app)
	secretsRepoMock.AssertExpectations(t)

	checkResponseStatus(t, resp, http.StatusInternalServerError)
}

func TestStoreAppStoreFault(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	app := roll.Application{
		ApplicationName: "ambivilant birds",
		DeveloperEmail:  "doug@dev.com",
		ClientID:        "1111-2222-3333333-4444444",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("StoreApplication", &app).Return(errors.New("storage fault"))

	secretsRepoMock := coreConfig.SecretsRepo.(*mocks.SecretsRepo)
	secretsRepoMock.On("StoreKeysForApp",
		mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()

	resp := testHTTPPut(t, addr+"/v1/applications/1111-2222-3333333-4444444", app)
	appRepoMock.AssertCalled(t, "StoreApplication", &app)
	secretsRepoMock.AssertExpectations(t)

	checkResponseStatus(t, resp, http.StatusInternalServerError)
}

func TestStoreAppNoResource(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	app := roll.Application{
		ApplicationName: "ambivilant birds",
		DeveloperEmail:  "doug@dev.com",
		ClientID:        "1111-2222-3333333-4444444",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	bodyReader := new(bytes.Buffer)
	enc := json.NewEncoder(bodyReader)
	err := enc.Encode(app)
	checkFatal(t, err)

	req, err := http.NewRequest("PUT", addr+"/v1/applications/", bodyReader)
	checkFatal(t, err)

	client := http.DefaultClient
	resp, err := client.Do(req)
	assert.Nil(t, err)

	checkResponseStatus(t, resp, http.StatusBadRequest)

}

func TestStoreAppBodyParseError(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	buf := new(bytes.Buffer)
	buf.WriteString(`{"this won't parse`)

	req, err := http.NewRequest("PUT", addr+"/v1/applications/", buf)
	checkFatal(t, err)

	client := http.DefaultClient
	resp, err := client.Do(req)
	assert.Nil(t, err)

	checkResponseStatus(t, resp, http.StatusBadRequest)

}

func TestAppUnsupportedMethod(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.Post(addr+"/v1/applications/1111-2222-3333333-4444444", "application/json", nil)
	assert.Nil(t, err)
	checkResponseStatus(t, resp, http.StatusMethodNotAllowed)
}

func TestStoreAppInvalidContent(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	app := roll.Application{
		ApplicationName: "ambivilant birds<script>",
		DeveloperEmail:  "doug@dev.com",
		ClientID:        "1111-2222-3333333-4444444",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	resp := testHTTPPut(t, addr+"/v1/applications/1111-2222-3333333-4444444", app)
	checkResponseStatus(t, resp, http.StatusBadRequest)
}

func TestGetApplication(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	returnVal := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		ClientID:        "1111-2222-3333333-4444444",
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
	assert.Equal(t, "1111-2222-3333333-4444444", actual.ClientID)
	assert.Equal(t, "fight club", actual.ApplicationName)
	assert.Equal(t, "not for browser clients", actual.ClientSecret)
	assert.Equal(t, "http://localhost:3000/ab", actual.RedirectURI)
	assert.Equal(t, "xtrac://localhost:9000", actual.LoginProvider)

}

func TestGetApplications(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	returnVal := []roll.Application{
		roll.Application{
			DeveloperEmail:  "doug@dev.com",
			ClientID:        "1111-2222-3333333-4444444",
			ApplicationName: "fight club",
			ClientSecret:    "not for browser clients",
			RedirectURI:     "http://localhost:3000/ab",
			LoginProvider:   "xtrac://localhost:9000",
		},
		roll.Application{
			DeveloperEmail:  "doug@dev.com",
			ClientID:        "1111-2222-3333333-4444444",
			ApplicationName: "fight club",
			ClientSecret:    "not for browser clients",
			RedirectURI:     "http://localhost:3000/ab",
			LoginProvider:   "xtrac://localhost:9000",
		}}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("ListApplications").Return(returnVal, nil)

	resp := testHTTPGet(t, addr+"/v1/applications/", nil)
	appRepoMock.AssertCalled(t, "ListApplications")

	var actual []roll.Application

	checkResponseBody(t, resp, &actual)
	for _, app := range actual {
		assert.Equal(t, "doug@dev.com", app.DeveloperEmail)
		assert.Equal(t, "1111-2222-3333333-4444444", app.ClientID)
		assert.Equal(t, "fight club", app.ApplicationName)
		assert.Equal(t, "not for browser clients", app.ClientSecret)
		assert.Equal(t, "http://localhost:3000/ab", app.RedirectURI)
		assert.Equal(t, "xtrac://localhost:9000", app.LoginProvider)
	}

}

func TestGetApplicationsReposError(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("ListApplications").Return(nil, errors.New("db error"))

	resp := testHTTPGet(t, addr+"/v1/applications/", nil)
	appRepoMock.AssertCalled(t, "ListApplications")

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

}

func TestRetrieveOfNonexistantApp(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(nil, nil)

	resp := testHTTPGet(t, addr+"/v1/applications/1111-2222-3333333-4444444", nil)
	appRepoMock.AssertCalled(t, "RetrieveApplication", "1111-2222-3333333-4444444")

	checkResponseStatus(t, resp, http.StatusNotFound)
}

func TestErrorOnAppRetrieve(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(nil, errors.New("big problem"))

	resp := testHTTPGet(t, addr+"/v1/applications/1111-2222-3333333-4444444", nil)
	appRepoMock.AssertCalled(t, "RetrieveApplication", "1111-2222-3333333-4444444")

	checkResponseStatus(t, resp, http.StatusInternalServerError)
}
