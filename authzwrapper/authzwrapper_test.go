package authzwrapper

import (
	"github.com/stretchr/testify/assert"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/roll/roll/mocks"
	"github.com/xtraclabs/rollsecrets/secrets"
	rolltoken "github.com/xtraclabs/rollsecrets/token"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func echoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body.Close()
			w.Write(body)
			w.Write([]byte("\n"))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func TestNoToken(t *testing.T) {
	secretsRepo := new(mocks.SecretsRepo)
	adminRepo := new(mocks.AdminRepo)
	testServer := httptest.NewServer(Wrap(secretsRepo, adminRepo, []string{}, echoHandler()))
	defer testServer.Close()

	resp, err := http.Post(testServer.URL, "text/plain", nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestGoodToken(t *testing.T) {

	app := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		ClientID:        "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		ClientSecret:    "not for browser clients",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	appRepoMock := new(mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(&app, nil)

	privateKey, publicKey, err := secrets.GenerateKeyPair()
	assert.Nil(t, err)

	secretsMock := new(mocks.SecretsRepo)
	secretsMock.On("RetrievePrivateKeyForApp", "1111-2222-3333333-4444444").Return(privateKey, nil)
	secretsMock.On("RetrievePublicKeyForApp", "1111-2222-3333333-4444444").Return(publicKey, nil)

	adminRepo := new(mocks.AdminRepo)

	token, err := rolltoken.GenerateToken("a-subject", "", app.ClientID, app.ApplicationName, privateKey)
	assert.Nil(t, err)

	testServer := httptest.NewServer(Wrap(secretsMock, adminRepo, []string{}, echoHandler()))
	defer testServer.Close()

	client := http.Client{}
	req, err := http.NewRequest("POST", testServer.URL, nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

}

func TestUnsecureWithSubjectHeader(t *testing.T) {
	testServer := httptest.NewServer(WrapUnsecure(echoHandler()))
	defer testServer.Close()

	client := http.Client{}
	req, err := http.NewRequest("POST", testServer.URL, nil)
	assert.Nil(t, err)
	req.Header.Add("X-Roll-Subject", "foobar")

	resp, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUnsecureMissingSubjectHeader(t *testing.T) {
	testServer := httptest.NewServer(WrapUnsecure(echoHandler()))
	defer testServer.Close()

	client := http.Client{}
	req, err := http.NewRequest("POST", testServer.URL, nil)
	assert.Nil(t, err)

	resp, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestMalformedToken(t *testing.T) {
	secretsRepo := new(mocks.SecretsRepo)
	adminRepo := new(mocks.AdminRepo)

	testServer := httptest.NewServer(Wrap(secretsRepo, adminRepo, []string{}, echoHandler()))
	defer testServer.Close()

	client := http.Client{}
	req, err := http.NewRequest("POST", testServer.URL, nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Bearer "+"nonsense-and-not-a-JWT")

	resp, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestNonBearerToken(t *testing.T) {
	secretsRepo := new(mocks.SecretsRepo)
	adminRepo := new(mocks.AdminRepo)
	testServer := httptest.NewServer(Wrap(secretsRepo, adminRepo, []string{}, echoHandler()))
	defer testServer.Close()

	client := http.Client{}
	req, err := http.NewRequest("POST", testServer.URL, nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "nonsense-and-not-a-JWT")

	resp, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestInvalidSignature(t *testing.T) {
	app := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		ClientID:        "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		ClientSecret:    "not for browser clients",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	appRepoMock := new(mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(&app, nil)

	privateKey, publicKey, err := secrets.GenerateKeyPair()
	assert.Nil(t, err)

	private2, _, err := secrets.GenerateKeyPair()
	assert.Nil(t, err)

	secretsMock := new(mocks.SecretsRepo)
	secretsMock.On("RetrievePrivateKeyForApp", "1111-2222-3333333-4444444").Return(privateKey, nil)
	secretsMock.On("RetrievePublicKeyForApp", "1111-2222-3333333-4444444").Return(publicKey, nil)

	adminRepo := new(mocks.AdminRepo)

	token, err := rolltoken.GenerateToken("b-subject", "", app.ClientID, app.ApplicationName, private2)
	assert.Nil(t, err)

	testServer := httptest.NewServer(Wrap(secretsMock, adminRepo, []string{}, echoHandler()))
	defer testServer.Close()

	client := http.Client{}
	req, err := http.NewRequest("POST", testServer.URL, nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuthCodeUsedForAccess(t *testing.T) {
	app := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		ClientID:        "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		ClientSecret:    "not for browser clients",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	appRepoMock := new(mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(&app, nil)

	privateKey, publicKey, err := secrets.GenerateKeyPair()
	assert.Nil(t, err)

	secretsMock := new(mocks.SecretsRepo)
	secretsMock.On("RetrievePrivateKeyForApp", "1111-2222-3333333-4444444").Return(privateKey, nil)
	secretsMock.On("RetrievePublicKeyForApp", "1111-2222-3333333-4444444").Return(publicKey, nil)

	adminRepo := new(mocks.AdminRepo)

	token, err := rolltoken.GenerateCode("a-subject", "", app.ClientID, privateKey)
	assert.Nil(t, err)

	testServer := httptest.NewServer(Wrap(secretsMock, adminRepo, []string{}, echoHandler()))
	defer testServer.Close()

	client := http.Client{}
	req, err := http.NewRequest("POST", testServer.URL, nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
