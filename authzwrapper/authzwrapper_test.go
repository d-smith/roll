package authzwrapper

import (
	"github.com/stretchr/testify/assert"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/roll/roll/mocks"
	"github.com/xtraclabs/roll/secrets"
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
	testServer := httptest.NewServer(Wrap(echoHandler(), secretsRepo))
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

	token, err := roll.GenerateToken(&app, privateKey)
	assert.Nil(t, err)

	testServer := httptest.NewServer(Wrap(echoHandler(), secretsMock))
	defer testServer.Close()

	client := http.Client{}
	req, err := http.NewRequest("POST", testServer.URL, nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

}

func TestInvalidToken(t *testing.T) {
	secretsRepo := new(mocks.SecretsRepo)
	testServer := httptest.NewServer(Wrap(echoHandler(), secretsRepo))
	defer testServer.Close()

	client := http.Client{}
	req, err := http.NewRequest("POST", testServer.URL, nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Bearer "+"nonsense-and-not-a-JWT")

	resp, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestExpiredToken(t *testing.T) {

}

func TestInvalidTokenSignature(t *testing.T) {

}
