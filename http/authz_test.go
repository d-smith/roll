package http

import (
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/roll/roll/mocks"
	"github.com/xtraclabs/roll/secrets"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRequiredQueryParamsPresent(t *testing.T) {
	t.Log("given a request with none of the required query params")
	t.Log("when requiredQueryParamsPresent is called")
	t.Log("then false is returned")

	req, _ := http.NewRequest("GET", "/", nil)
	assert.False(t, requiredQueryParamsPresent(req))

	t.Log("given a request with some but not all of the required query params")
	t.Log("when requiredQueryParamsPresent is called")
	t.Log("then false is returned")
	req, _ = http.NewRequest("GET", "/?client_id=123", nil)
	assert.False(t, requiredQueryParamsPresent(req))

	t.Log("given a request with all of the required query params")
	t.Log("when requiredQueryParamsPresent is called")
	t.Log("then true is returned")
	req, _ = http.NewRequest("POST", "/?client_id=123&redirect_uri=x&response_type=X", nil)
	assert.True(t, requiredQueryParamsPresent(req))
}

func TestInputParamsValid(t *testing.T) {
	core, coreConfig := NewTestCore()

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

	req, _ := http.NewRequest("POST", "/?client_id=1111-2222-3333333-4444444&redirect_uri=http://localhost:3000/ab&response_type=token", nil)

	app, err := validateInputParams(core, req)
	assert.Nil(t, err)
	assert.NotNil(t, app)
}

func TestInputParamsInvalidResponseType(t *testing.T) {
	req, _ := http.NewRequest("POST", "/?client_id=1111-2222-3333333-4444444&redirect_uri=http://localhost:3000/ab&response_type=bad", nil)
	app, err := validateInputParams(nil, req)
	assert.Nil(t, app)
	assert.NotNil(t, err)
}

func TestInputParamsMissingResponseType(t *testing.T) {
	req, _ := http.NewRequest("POST", "/?client_id=1111-2222-3333333-4444444&redirect_uri=http://localhost:3000/ab", nil)
	app, err := validateInputParams(nil, req)
	assert.Nil(t, app)
	assert.NotNil(t, err)
}

func TestInputParamsMissingClientID(t *testing.T) {
	core, coreConfig := NewTestCore()
	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "").Return(nil, nil)

	req, _ := http.NewRequest("POST", "/?redirect_uri=http://localhost:3000/ab&response_type=code", nil)
	app, err := validateInputParams(core, req)
	assert.Nil(t, app)
	assert.NotNil(t, err)
	println(err.Error())
}

func TestInputParamsNoSuchClientId(t *testing.T) {
	core, coreConfig := NewTestCore()

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(nil, errors.New("whoops"))

	req, _ := http.NewRequest("POST", "/?client_id=1111-2222-3333333-4444444&redirect_uri=http://localhost:3000/ab&response_type=token", nil)

	app, err := validateInputParams(core, req)
	assert.NotNil(t, err)
	assert.Nil(t, app)
}

func TestInputParamsInvalidRedirectURI(t *testing.T) {
	core, coreConfig := NewTestCore()

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

	req, _ := http.NewRequest("POST", "/?client_id=1111-2222-3333333-4444444&redirect_uri=bogus&response_type=token", nil)
	app, err := validateInputParams(core, req)
	assert.NotNil(t, err)
	assert.Nil(t, app)

}

func TestExecuteAuthTemplateForCode(t *testing.T) {
	w := httptest.NewRecorder()
	pageCtx := &authPageContext{
		AppName:  "test-application-name",
		ClientID: "test-application-key",
	}

	req, _ := http.NewRequest("POST", "/?client_id=1111-2222-3333333-4444444&redirect_uri=bogus&response_type=code", nil)

	err := executeAuthTemplate(w, req, pageCtx)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.True(t, strings.Contains(body, `name="client_id" value="test-application-key"`))
	assert.True(t, strings.Contains(body, ` <h2>test-application-name`))
	assert.True(t, strings.Contains(body, `name="response_type" value="code"`))
}

func TestExecuteAuthTemplateMissingResponseType(t *testing.T) {
	w := httptest.NewRecorder()
	pageCtx := &authPageContext{
		AppName:  "test-application-name",
		ClientID: "test-application-key",
	}

	req, _ := http.NewRequest("POST", "/?client_id=1111-2222-3333333-4444444&redirect_uri=bogus", nil)

	err := executeAuthTemplate(w, req, pageCtx)
	assert.NotNil(t, err)

}

func TestExecuteAuthTemplateBogusResponseType(t *testing.T) {
	w := httptest.NewRecorder()
	pageCtx := &authPageContext{
		AppName:  "test-application-name",
		ClientID: "test-application-key",
	}

	req, _ := http.NewRequest("POST", "/?client_id=1111-2222-3333333-4444444&redirect_uri=bogus&response_type=bogus", nil)

	err := executeAuthTemplate(w, req, pageCtx)
	assert.NotNil(t, err)

}

func TestHandleAuthorizeUnsupportedMethod(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.Post(addr+"/oauth2/authorize", "", nil)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestHandleAuthorizeMissingParams(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp := TestHTTPGet(t, addr+"/oauth2/authorize", nil)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAuthValidateMissingParams(t *testing.T) {
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

	resp, err := http.PostForm(addr+"/oauth2/validate",
		url.Values{"username": {},
			"password": {}})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "Expected single response_type param as part of query params"))
}

func TestAuthValidateBadResponseType(t *testing.T) {
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

	resp, err := http.PostForm(addr+"/oauth2/validate",
		url.Values{"username": {"x"},
			"password":      {"y"},
			"authorize":     {"allow"},
			"response_type": {"bad"},
			"client_id":     {"111-22-33"}})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "valid values for response_type are token and code"))
}

func TestAuthValidateCodeResponseAuthenticateOk(t *testing.T) {

	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	privateKey, publicKey, err := secrets.GenerateKeyPair()
	assert.Nil(t, err)

	secretsMock := coreConfig.SecretsRepo.(*mocks.SecretsRepo)
	secretsMock.On("RetrievePrivateKeyForApp", "1111-2222-3333333-4444444").Return(privateKey, nil)
	secretsMock.On("RetrievePublicKeyForApp", "1111-2222-3333333-4444444").Return(publicKey, nil)

	var loginCalled = false
	ls := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loginCalled = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ls.Close()

	//TODO - use a second callback where we serve up a script to extract the page details sent
	//on deny and post those details to another test server.
	var callbackInvoked = false

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callbackInvoked = true
		code := r.FormValue("code")
		token, err := jwt.Parse(code, roll.GenerateKeyExtractionFunction(core.SecretsRepo))
		assert.Nil(t, err)
		scope, ok := token.Claims["scope"].(string)
		assert.True(t, ok)
		assert.Equal(t, "xtAuthCode", scope)
	}))
	defer ts.Close()

	lsURL, _ := url.Parse(ls.URL)

	returnVal := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		ClientID:        "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		ClientSecret:    "not for browser clients",
		RedirectURI:     ts.URL,
		LoginProvider:   "xtrac://" + lsURL.Host,
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	_, err = http.PostForm(addr+"/oauth2/validate",
		url.Values{"username": {"x"},
			"password":      {"y"},
			"authorize":     {"allow"},
			"response_type": {"code"},
			"client_id":     {"1111-2222-3333333-4444444"}})
	assert.Nil(t, err)
	assert.True(t, callbackInvoked)
	assert.True(t, loginCalled)
}

func TestAuthValidateCodeResponseAuthenticateAdminScopeOk(t *testing.T) {

	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	privateKey, publicKey, err := secrets.GenerateKeyPair()
	assert.Nil(t, err)

	secretsMock := coreConfig.SecretsRepo.(*mocks.SecretsRepo)
	secretsMock.On("RetrievePrivateKeyForApp", "1111-2222-3333333-4444444").Return(privateKey, nil)
	secretsMock.On("RetrievePublicKeyForApp", "1111-2222-3333333-4444444").Return(publicKey, nil)

	var loginCalled = false
	ls := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loginCalled = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ls.Close()

	var callbackInvoked = false

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callbackInvoked = true
		code := r.FormValue("code")
		token, err := jwt.Parse(code, roll.GenerateKeyExtractionFunction(core.SecretsRepo))
		assert.Nil(t, err)
		scope, ok := token.Claims["scope"].(string)
		assert.True(t, ok)
		assert.Equal(t, "xtAuthCode admin", scope)
	}))
	defer ts.Close()

	lsURL, _ := url.Parse(ls.URL)

	returnVal := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		ClientID:        "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		ClientSecret:    "not for browser clients",
		RedirectURI:     ts.URL,
		LoginProvider:   "xtrac://" + lsURL.Host,
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	adminRepoMock := coreConfig.AdminRepo.(*mocks.AdminRepo)
	adminRepoMock.On("IsAdmin", "x").Return(true, nil)

	_, err = http.PostForm(addr+"/oauth2/validate",
		url.Values{"username": {"x"},
			"password":      {"y"},
			"authorize":     {"allow"},
			"response_type": {"code"},
			"scope":         {"admin"},
			"client_id":     {"1111-2222-3333333-4444444"}})
	assert.Nil(t, err)
	assert.True(t, callbackInvoked)
	assert.True(t, loginCalled)
}

func TestAuthValidateCodeResponseAuthenticateAdminScopeDenied(t *testing.T) {

	var loginCalled = false
	ls := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loginCalled = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ls.Close()

	//TODO - use a second callback where we serve up a script to extract the page details sent
	//on deny and post those details to another test server.
	var callbackInvoked = false

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callbackInvoked = true
		assert.True(t, strings.Contains(r.RequestURI, "error=access_denied&error_description=scope-problem"))
	}))
	defer ts.Close()

	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	lsURL, _ := url.Parse(ls.URL)

	returnVal := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		ClientID:        "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		ClientSecret:    "not for browser clients",
		RedirectURI:     ts.URL + "/foo",
		LoginProvider:   "xtrac://" + lsURL.Host,
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	adminRepoMock := coreConfig.AdminRepo.(*mocks.AdminRepo)
	adminRepoMock.On("IsAdmin", "x").Return(false, nil)

	privateKey, _, err := secrets.GenerateKeyPair()
	assert.Nil(t, err)

	secretsMock := coreConfig.SecretsRepo.(*mocks.SecretsRepo)
	secretsMock.On("RetrievePrivateKeyForApp", "1111-2222-3333333-4444444").Return(privateKey, nil)

	_, err = http.PostForm(addr+"/oauth2/validate",
		url.Values{"username": {"x"},
			"password":      {"y"},
			"authorize":     {"allow"},
			"response_type": {"code"},
			"scope":         {"admin"},
			"client_id":     {"1111-2222-3333333-4444444"}})
	assert.Nil(t, err)
	assert.True(t, callbackInvoked)
}

func TestAuthValidateCodeResponseAuthenticateAdminScopeError(t *testing.T) {

	var loginCalled = false
	ls := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loginCalled = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ls.Close()

	//TODO - use a second callback where we serve up a script to extract the page details sent
	//on deny and post those details to another test server.
	var callbackInvoked = false

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callbackInvoked = true
		assert.True(t, strings.Contains(r.RequestURI, "error=server_error&error_description="))
	}))
	defer ts.Close()

	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	lsURL, _ := url.Parse(ls.URL)

	returnVal := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		ClientID:        "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		ClientSecret:    "not for browser clients",
		RedirectURI:     ts.URL + "/foo",
		LoginProvider:   "xtrac://" + lsURL.Host,
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	adminRepoMock := coreConfig.AdminRepo.(*mocks.AdminRepo)
	adminRepoMock.On("IsAdmin", "x").Return(false, errors.New("BOOM!"))

	privateKey, _, err := secrets.GenerateKeyPair()
	assert.Nil(t, err)

	secretsMock := coreConfig.SecretsRepo.(*mocks.SecretsRepo)
	secretsMock.On("RetrievePrivateKeyForApp", "1111-2222-3333333-4444444").Return(privateKey, nil)

	_, err = http.PostForm(addr+"/oauth2/validate",
		url.Values{"username": {"x"},
			"password":      {"y"},
			"authorize":     {"allow"},
			"response_type": {"code"},
			"scope":         {"admin"},
			"client_id":     {"1111-2222-3333333-4444444"}})
	assert.Nil(t, err)
	assert.True(t, callbackInvoked)
}
