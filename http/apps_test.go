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

func TestStoreAppOK(t *testing.T) {
	t.Log("TestStoreAppOK")
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	app := roll.Application{
		ApplicationName: "ambivilant birds",
		ClientID:        "steve",
		DeveloperEmail:  "doug@dev.com",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("CreateApplication", &app).Return(nil)

	secretsRepoMock := coreConfig.SecretsRepo.(*mocks.SecretsRepo)
	secretsRepoMock.On("StoreKeysForApp",
		mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()

	resp := TestHTTPPost(t, addr+"/v1/applications", app)
	appRepoMock.AssertCalled(t, "CreateApplication", &app)
	appRepoMock.AssertExpectations(t)
	secretsRepoMock.AssertExpectations(t)

	checkResponseStatus(t, resp, http.StatusOK)

	var cid clientID

	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	err := dec.Decode(&cid)
	assert.Nil(t, err)
	assert.Equal(t, "steve", cid.ClientID)

}

func TestUpdateAppOK(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	app := roll.Application{
		ApplicationName: "ambivilant birds",
		ClientID:        "111-222-333",
		DeveloperEmail:  "doug@dev.com",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	app2 := roll.Application{
		ApplicationName: "foos",
		ClientID:        "111-222-333",
		DeveloperEmail:  "doug@dev.com",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "111-222-333").Return(&app, nil)
	appRepoMock.On("UpdateApplication", &app2).Return(nil)

	secretsRepoMock := coreConfig.SecretsRepo.(*mocks.SecretsRepo)
	secretsRepoMock.On("StoreKeysForApp",
		mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()

	resp := TestHTTPPut(t, addr+"/v1/applications/111-222-333", app2)
	appRepoMock.AssertCalled(t, "RetrieveApplication", "111-222-333")
	appRepoMock.AssertCalled(t, "UpdateApplication", &app2)
	secretsRepoMock.AssertNotCalled(t, "StoreKeysForApp")

	checkResponseStatus(t, resp, http.StatusNoContent)
}

func TestUpdateAppStoreFault(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	app := roll.Application{
		ApplicationName: "ambivilant birds",
		ClientID:        "111-222-333",
		DeveloperEmail:  "doug@dev.com",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	app2 := roll.Application{
		ApplicationName: "foos",
		ClientID:        "111-222-333",
		DeveloperEmail:  "doug@dev.com",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "111-222-333").Return(&app, nil)
	appRepoMock.On("UpdateApplication", &app2).Return(errors.New("boom!"))

	secretsRepoMock := coreConfig.SecretsRepo.(*mocks.SecretsRepo)
	secretsRepoMock.On("StoreKeysForApp",
		mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()

	resp := TestHTTPPut(t, addr+"/v1/applications/111-222-333", app2)
	appRepoMock.AssertCalled(t, "RetrieveApplication", "111-222-333")
	appRepoMock.AssertCalled(t, "UpdateApplication", &app2)
	secretsRepoMock.AssertNotCalled(t, "StoreKeysForApp")

	checkResponseStatus(t, resp, http.StatusInternalServerError)
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

	resp := TestHTTPPost(t, addr+"/v1/applications", app)
	secretsRepoMock.AssertExpectations(t)

	checkResponseStatus(t, resp, http.StatusInternalServerError)
}

func TestUpdateAppRetrieveFault(t *testing.T) {
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
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(nil, errors.New("kaboom"))

	resp := TestHTTPPut(t, addr+"/v1/applications/1111-2222-3333333-4444444", app)

	checkResponseStatus(t, resp, http.StatusInternalServerError)
}

func TestUpdateAppNotFound(t *testing.T) {
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
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(nil, nil)

	resp := TestHTTPPut(t, addr+"/v1/applications/1111-2222-3333333-4444444", app)

	checkResponseStatus(t, resp, http.StatusNotFound)
}

type BadIDGenerator struct{}

func (big BadIDGenerator) GenerateID() (string, error) {
	return "", errors.New("whoops")
}

func TestStoreAppIDGenFault(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	core.IdGenerator = BadIDGenerator{}

	app := roll.Application{
		ApplicationName: "ambivilant birds",
		DeveloperEmail:  "doug@dev.com",
		ClientID:        "1111-2222-3333333-4444444",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	resp := TestHTTPPost(t, addr+"/v1/applications", app)

	checkResponseStatus(t, resp, http.StatusInternalServerError)
}

func TestStoreAppContentValidationFault(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	app := roll.Application{
		ApplicationName: "ambivilant birds",
		DeveloperEmail:  "yeah yeah yeah",
		ClientID:        "1111-2222-3333333-4444444",
		RedirectURI:     "not a uri",
		LoginProvider:   "xtrac://localhost:9000",
	}

	resp := TestHTTPPost(t, addr+"/v1/applications", app)

	checkResponseStatus(t, resp, http.StatusBadRequest)
}

func TestStoreAppStoreFault(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	app := roll.Application{
		ApplicationName: "ambivilant birds",
		DeveloperEmail:  "doug@dev.com",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("CreateApplication", mock.AnythingOfType("*roll.Application")).Return(errors.New("storage fault"))

	secretsRepoMock := coreConfig.SecretsRepo.(*mocks.SecretsRepo)
	secretsRepoMock.On("StoreKeysForApp",
		mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()

	resp := TestHTTPPost(t, addr+"/v1/applications", app)
	appRepoMock.AssertCalled(t, "CreateApplication", mock.AnythingOfType("*roll.Application"))
	secretsRepoMock.AssertExpectations(t)

	checkResponseStatus(t, resp, http.StatusInternalServerError)
}

func TestUpdateAppNoResource(t *testing.T) {
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

	req, err := http.NewRequest("POST", addr+"/v1/applications", buf)
	checkFatal(t, err)

	client := http.DefaultClient
	resp, err := client.Do(req)
	assert.Nil(t, err)

	checkResponseStatus(t, resp, http.StatusBadRequest)

}

func TestUpdateAppBodyParseError(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	buf := new(bytes.Buffer)
	buf.WriteString(`{"this won't parse`)

	req, err := http.NewRequest("PUT", addr+"/v1/applications/11-22-33", buf)
	checkFatal(t, err)

	client := http.DefaultClient
	resp, err := client.Do(req)
	assert.Nil(t, err)

	checkResponseStatus(t, resp, http.StatusBadRequest)

}

func TestAppUnsupportedMethodBaseURI(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	req, err := http.NewRequest("OPTIONS", addr+"/v1/applications", nil)
	assert.Nil(t, err)

	client := http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	checkResponseStatus(t, resp, http.StatusMethodNotAllowed)
}

func TestAppUnsupportedMethodExtendedURI(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	req, err := http.NewRequest("OPTIONS", addr+"/v1/applications/", nil)
	assert.Nil(t, err)

	client := http.Client{}
	resp, err := client.Do(req)
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

	resp := TestHTTPPut(t, addr+"/v1/applications/1111-2222-3333333-4444444", app)
	checkResponseStatus(t, resp, http.StatusBadRequest)
}

func TestGetApplication(t *testing.T) {
	t.Log("xxxxxxxxxxxxxxxxxxxxxxx")
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

	t.Log("set up mock app repo")
	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	t.Log("get get get get get")
	resp := TestHTTPGet(t, addr+"/v1/applications/1111-2222-3333333-4444444", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	t.Log("assert get was called with the input client id")
	appRepoMock.AssertCalled(t, "RetrieveApplication", "1111-2222-3333333-4444444")

	var actual roll.Application

	t.Log("check the response body")
	checkResponseBody(t, resp, &actual)
	assert.Equal(t, "doug@dev.com", actual.DeveloperEmail)
	assert.Equal(t, "1111-2222-3333333-4444444", actual.ClientID)
	assert.Equal(t, "fight club", actual.ApplicationName)
	assert.Equal(t, "not for browser clients", actual.ClientSecret)
	assert.Equal(t, "http://localhost:3000/ab", actual.RedirectURI)
	assert.Equal(t, "xtrac://localhost:9000", actual.LoginProvider)

}

func TestGetApplicationsOK(t *testing.T) {
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

	resp := TestHTTPGet(t, addr+"/v1/applications", nil)
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

	resp := TestHTTPGet(t, addr+"/v1/applications", nil)
	appRepoMock.AssertCalled(t, "ListApplications")

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

}

func TestRetrieveOfNonexistantApp(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(nil, nil)

	resp := TestHTTPGet(t, addr+"/v1/applications/1111-2222-3333333-4444444", nil)
	appRepoMock.AssertCalled(t, "RetrieveApplication", "1111-2222-3333333-4444444")

	checkResponseStatus(t, resp, http.StatusNotFound)
}

func TestErrorOnAppRetrieve(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(nil, errors.New("big problem"))

	resp := TestHTTPGet(t, addr+"/v1/applications/1111-2222-3333333-4444444", nil)
	appRepoMock.AssertCalled(t, "RetrieveApplication", "1111-2222-3333333-4444444")

	checkResponseStatus(t, resp, http.StatusInternalServerError)
}
