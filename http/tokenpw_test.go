package http

import (
	"encoding/json"
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/roll/roll/mocks"
	"github.com/xtraclabs/rollsecrets/secrets"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestPWGrantMissingClientID(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"password"}})

	assert.Nil(t, err)
	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "client_id missing from request"))
}

func TestPWGrantMissingClientSecretX(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"password"},
			"client_id": {"1"}})

	assert.Nil(t, err)
	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "client_secret missing from request"))
}

func TestPWGrantMissingUsername(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"password"},
			"client_id":     {"1"},
			"client_secret": {"2"}})

	assert.Nil(t, err)
	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "username missing from request"))
}

func TestPWGrantMissingPassword(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"password"},
			"client_id":     {"1"},
			"client_secret": {"2"},
			"username":      {"foo"}})

	assert.Nil(t, err)
	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "password missing from request"))
}

func TestTokenInvalidMethod(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.Get(addr + OAuth2TokenBaseURI)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestPWGrantLoginFailure(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	var loginCalled = false
	ls := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loginCalled = true
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ls.Close()

	lsURL, _ := url.Parse(ls.URL)

	returnVal := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		ClientID:        "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		ClientSecret:    "not for browser clients",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://" + lsURL.Host,
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"password"},
			"client_id":     {"1111-2222-3333333-4444444"},
			"client_secret": {"not for browser clients"},
			"username":      {"abc"},
			"password":      {"xxxxxxxx"}})

	assert.True(t, loginCalled)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestPWGrantLoginCallError(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	returnVal := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		ClientID:        "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		ClientSecret:    "not for browser clients",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:77",
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"password"},
			"client_id":     {"1111-2222-3333333-4444444"},
			"client_secret": {"not for browser clients"},
			"username":      {"abc"},
			"password":      {"xxxxxxxx"}})

	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestPWGrantLoginOk(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	var loginCalled = false
	ls := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loginCalled = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ls.Close()

	lsURL, _ := url.Parse(ls.URL)

	returnVal := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		ClientID:        "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		ClientSecret:    "not for browser clients",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://" + lsURL.Host,
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	privateKey, publicKey, err := secrets.GenerateKeyPair()
	assert.Nil(t, err)

	secretsMock := coreConfig.SecretsRepo.(*mocks.SecretsRepo)
	secretsMock.On("RetrievePrivateKeyForApp", "1111-2222-3333333-4444444").Return(privateKey, nil)
	secretsMock.On("RetrievePublicKeyForApp", "1111-2222-3333333-4444444").Return(publicKey, nil)

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"password"},
			"client_id":     {"1111-2222-3333333-4444444"},
			"client_secret": {"not for browser clients"},
			"username":      {"abc"},
			"password":      {"xxxxxxxx"}})

	assert.True(t, loginCalled)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body := responseAsString(t, resp)
	var jsonResponse accessTokenResponse
	err = json.Unmarshal([]byte(body), &jsonResponse)
	assert.Nil(t, err)
	assert.True(t, jsonResponse.AccessToken != "")
	assert.True(t, jsonResponse.TokenType == "Bearer")

	token, err := jwt.Parse(jsonResponse.AccessToken, roll.GenerateKeyExtractionFunction(core.SecretsRepo))
	assert.Nil(t, err)
	fmt.Println(token.Claims)
	assert.Equal(t, "1111-2222-3333333-4444444", token.Claims["aud"].(string))
	assert.Equal(t, "abc", token.Claims["sub"].(string))
	assert.Equal(t, "", token.Claims["scope"].(string))
}

func TestPWGrantLoginOkAdminScope(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	var loginCalled = false
	ls := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loginCalled = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ls.Close()

	lsURL, _ := url.Parse(ls.URL)

	returnVal := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		ClientID:        "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		ClientSecret:    "not for browser clients",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://" + lsURL.Host,
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	privateKey, publicKey, err := secrets.GenerateKeyPair()
	assert.Nil(t, err)

	secretsMock := coreConfig.SecretsRepo.(*mocks.SecretsRepo)
	secretsMock.On("RetrievePrivateKeyForApp", "1111-2222-3333333-4444444").Return(privateKey, nil)
	secretsMock.On("RetrievePublicKeyForApp", "1111-2222-3333333-4444444").Return(publicKey, nil)

	adminRepoMock := coreConfig.AdminRepo.(*mocks.AdminRepo)
	adminRepoMock.On("IsAdmin", "abc").Return(true, nil)

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"password"},
			"client_id":     {"1111-2222-3333333-4444444"},
			"client_secret": {"not for browser clients"},
			"username":      {"abc"},
			"scope":         {"admin"},
			"password":      {"xxxxxxxx"}})

	assert.True(t, loginCalled)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body := responseAsString(t, resp)
	var jsonResponse accessTokenResponse
	err = json.Unmarshal([]byte(body), &jsonResponse)
	assert.Nil(t, err)
	assert.True(t, jsonResponse.AccessToken != "")
	assert.True(t, jsonResponse.TokenType == "Bearer")

	token, err := jwt.Parse(jsonResponse.AccessToken, roll.GenerateKeyExtractionFunction(core.SecretsRepo))
	assert.Nil(t, err)
	fmt.Println(token.Claims)
	assert.Equal(t, "1111-2222-3333333-4444444", token.Claims["aud"].(string))
	assert.Equal(t, "abc", token.Claims["sub"].(string))
	assert.Equal(t, "admin", token.Claims["scope"].(string))
}

func TestPWGrantLoginOkInvalidScope(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	var loginCalled = false
	ls := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loginCalled = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ls.Close()

	lsURL, _ := url.Parse(ls.URL)

	returnVal := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		ClientID:        "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		ClientSecret:    "not for browser clients",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://" + lsURL.Host,
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	privateKey, publicKey, err := secrets.GenerateKeyPair()
	assert.Nil(t, err)

	secretsMock := coreConfig.SecretsRepo.(*mocks.SecretsRepo)
	secretsMock.On("RetrievePrivateKeyForApp", "1111-2222-3333333-4444444").Return(privateKey, nil)
	secretsMock.On("RetrievePublicKeyForApp", "1111-2222-3333333-4444444").Return(publicKey, nil)

	adminRepoMock := coreConfig.AdminRepo.(*mocks.AdminRepo)
	adminRepoMock.On("IsAdmin", "abc").Return(true, nil)

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"password"},
			"client_id":     {"1111-2222-3333333-4444444"},
			"client_secret": {"not for browser clients"},
			"username":      {"abc"},
			"scope":         {"adminxxx"},
			"password":      {"xxxxxxxx"}})

	assert.True(t, loginCalled)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

}

func TestPWGrantAppLookupErr(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(nil, errors.New("Drat"))

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"password"},
			"client_id":     {"1111-2222-3333333-4444444"},
			"client_secret": {"not for browser clients"},
			"username":      {"luser"},
			"password":      {"xxxxxxxx"}})

	assert.Nil(t, err)
	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "Missing or invalid form data"))
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestPWGrantInvalidClientSecret(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	returnVal := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		ClientID:        "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		ClientSecret:    "hax0r",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"password"},
			"client_id":     {"1111-2222-3333333-4444444"},
			"client_secret": {"not for browser clients"},
			"username":      {"luser"},
			"password":      {"xxxxxxxx"}})

	assert.Nil(t, err)
	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "Invalid application details"))
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
