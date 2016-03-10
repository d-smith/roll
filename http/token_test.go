package http

import (
	"encoding/json"
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/roll/roll/mocks"
	"github.com/xtraclabs/rollsecrets/secrets"
	rolltoken "github.com/xtraclabs/rollsecrets/token"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestTokenMissingGrantType(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{})

	assert.Nil(t, err)
	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "Invalid grant_type"))
}

func TestTokenInvalidGrantType(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"foo grant"}})

	assert.Nil(t, err)
	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "Invalid grant_type"))
}

func TestTokenMissingClientID(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"authorization_code"}})

	assert.Nil(t, err)
	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "client_id missing from request"))
}

func TestTokenMissingClientSecretX(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"authorization_code"},
			"client_id": {"1"}})

	assert.Nil(t, err)
	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "client_secret missing from request"))
}

func TestTokenMissingRedirectUri(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"authorization_code"},
			"client_id":     {"1"},
			"client_secret": {"xxx"}})

	assert.Nil(t, err)
	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "redirect_uri missing from request"))
}

func TestTokenMissingCode(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"authorization_code"},
			"client_id":     {"1"},
			"client_secret": {"xxx"},
			"redirect_uri":  {"http://foo:1000"}})

	assert.Nil(t, err)
	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "code is missing from request"))
}

func TestTokenAppLookupErr(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(nil, errors.New("Drat"))

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"authorization_code"},
			"client_id":     {"1111-2222-3333333-4444444"},
			"client_secret": {"not for browser clients"},
			"redirect_uri":  {"http://localhost:3000/ab"},
			"code":          {"xxxxxxxx"}})

	assert.Nil(t, err)
	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "Missing or invalid form data"))
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestTokenAppLookupNil(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(nil, nil)

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"authorization_code"},
			"client_id":     {"1111-2222-3333333-4444444"},
			"client_secret": {"not for browser clients"},
			"redirect_uri":  {"http://localhost:3000/ab"},
			"code":          {"xxxxxxxx"}})

	assert.Nil(t, err)
	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "Invalid client id"))
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestTokenUnparsableAuthCode(t *testing.T) {
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
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"authorization_code"},
			"client_id":     {"1111-2222-3333333-4444444"},
			"client_secret": {"not for browser clients"},
			"redirect_uri":  {"http://localhost:3000/ab"},
			"code":          {"xxxxxxxx"}})

	assert.Nil(t, err)
	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "token contains an invalid number of segments"))
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestTokenInvalidClientSecret(t *testing.T) {
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
		url.Values{"grant_type": {"authorization_code"},
			"client_id":     {"1111-2222-3333333-4444444"},
			"client_secret": {"not for browser clients"},
			"redirect_uri":  {"http://localhost:3000/ab"},
			"code":          {"xxxxxxxx"}})

	assert.Nil(t, err)
	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "Invalid application details"))
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestTokenInvalidRedirect(t *testing.T) {
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
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"authorization_code"},
			"client_id":     {"1111-2222-3333333-4444444"},
			"client_secret": {"not for browser clients"},
			"redirect_uri":  {"xxxx"},
			"code":          {"xxxxxxxx"}})

	assert.Nil(t, err)
	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "Invalid application details"))
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestTokenSignedWithWrongKey(t *testing.T) {
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
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	privateKey, publicKey, err := secrets.GenerateKeyPair()
	assert.Nil(t, err)

	secretsMock := coreConfig.SecretsRepo.(*mocks.SecretsRepo)
	secretsMock.On("RetrievePrivateKeyForApp", "1111-2222-3333333-4444444").Return(privateKey, nil)
	secretsMock.On("RetrievePublicKeyForApp", "1111-2222-3333333-4444444").Return(publicKey, nil)

	otherKey, _, err := secrets.GenerateKeyPair()
	assert.Nil(t, err)

	code, err := rolltoken.GenerateCode("a-subject", "", returnVal.ClientID, otherKey)
	assert.Nil(t, err)

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"authorization_code"},
			"client_id":     {"1111-2222-3333333-4444444"},
			"client_secret": {"not for browser clients"},
			"redirect_uri":  {"http://localhost:3000/ab"},
			"code":          {code}})

	assert.Nil(t, err)
	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "verification error"))
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestTokenValidCode(t *testing.T) {
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
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	privateKey, publicKey, err := secrets.GenerateKeyPair()
	assert.Nil(t, err)

	secretsMock := coreConfig.SecretsRepo.(*mocks.SecretsRepo)
	secretsMock.On("RetrievePrivateKeyForApp", "1111-2222-3333333-4444444").Return(privateKey, nil)
	secretsMock.On("RetrievePublicKeyForApp", "1111-2222-3333333-4444444").Return(publicKey, nil)

	code, err := rolltoken.GenerateCode("b-subject", "", returnVal.ClientID, privateKey)
	assert.Nil(t, err)

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"authorization_code"},
			"client_id":     {"1111-2222-3333333-4444444"},
			"client_secret": {"not for browser clients"},
			"redirect_uri":  {"http://localhost:3000/ab"},
			"code":          {code}})

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body := responseAsString(t, resp)

	var jsonResponse accessTokenResponse
	err = json.Unmarshal([]byte(body), &jsonResponse)
	assert.Nil(t, err)
	assert.True(t, jsonResponse.AccessToken != "")
	assert.True(t, jsonResponse.TokenType == "Bearer")

}

func TestTokenValidCodeWithAdminScope(t *testing.T) {
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
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	privateKey, publicKey, err := secrets.GenerateKeyPair()
	assert.Nil(t, err)

	secretsMock := coreConfig.SecretsRepo.(*mocks.SecretsRepo)
	secretsMock.On("RetrievePrivateKeyForApp", "1111-2222-3333333-4444444").Return(privateKey, nil)
	secretsMock.On("RetrievePublicKeyForApp", "1111-2222-3333333-4444444").Return(publicKey, nil)

	code, err := rolltoken.GenerateCode("b-subject", "admin", returnVal.ClientID, privateKey)
	assert.Nil(t, err)

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"authorization_code"},
			"client_id":     {"1111-2222-3333333-4444444"},
			"client_secret": {"not for browser clients"},
			"redirect_uri":  {"http://localhost:3000/ab"},
			"code":          {code}})

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body := responseAsString(t, resp)

	var jsonResponse accessTokenResponse
	err = json.Unmarshal([]byte(body), &jsonResponse)
	assert.Nil(t, err)
	assert.True(t, jsonResponse.AccessToken != "")
	assert.True(t, jsonResponse.TokenType == "Bearer")

	token, err := jwt.Parse(jsonResponse.AccessToken, rolltoken.GenerateKeyExtractionFunction(core.SecretsRepo))
	assert.Nil(t, err)
	scope, ok := token.Claims["scope"].(string)
	assert.True(t, ok)
	assert.Equal(t, "admin", scope)

}
