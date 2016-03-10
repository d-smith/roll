package http

import (
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/roll/roll/mocks"
	"github.com/xtraclabs/rollsecrets/secrets"
	rolltoken "github.com/xtraclabs/rollsecrets/token"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestExecuteAuthTemplateForToken(t *testing.T) {
	w := httptest.NewRecorder()
	pageCtx := &authPageContext{
		AppName:  "test-application-name",
		ClientID: "test-application-key",
		Scope:    "scooby doo",
	}

	req, _ := http.NewRequest("POST", "/?client_id=1111-2222-3333333-4444444&redirect_uri=bogus&response_type=token&scope=admin", nil)

	err := executeAuthTemplate(w, req, pageCtx)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.True(t, strings.Contains(body, `name="client_id" value="test-application-key"`))
	assert.True(t, strings.Contains(body, ` <h2>test-application-name`))
	assert.True(t, strings.Contains(body, `name="response_type" value="token"`))
	assert.True(t, strings.Contains(body, `<input type="hidden" name="scope" value="scooby doo"`))
}

func TestHandleImpGrantAuthorize(t *testing.T) {
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

	resp := TestHTTPGet(t, addr+"/oauth2/authorize?client_id=1111-2222-3333333-4444444&redirect_uri=http://localhost:3000/ab&response_type=token", nil)
	appRepoMock.AssertCalled(t, "SystemRetrieveApplication", "1111-2222-3333333-4444444")

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bodyStr := responseAsString(t, resp)
	assert.True(t, strings.Contains(bodyStr, `name="client_id" value="1111-2222-3333333-4444444"`))
	assert.True(t, strings.Contains(bodyStr, ` <h2>fight club`))
	assert.True(t, strings.Contains(bodyStr, `name="response_type" value="token"`))
	assert.True(t, strings.Contains(bodyStr, `<input type="hidden" name="scope" value=""`))

}

func TestHandleImpGrantAuthorizeAdminScope(t *testing.T) {
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

	resp := TestHTTPGet(t, addr+"/oauth2/authorize?client_id=1111-2222-3333333-4444444&redirect_uri=http://localhost:3000/ab&response_type=token&scope=admin", nil)
	appRepoMock.AssertCalled(t, "SystemRetrieveApplication", "1111-2222-3333333-4444444")

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bodyStr := responseAsString(t, resp)
	assert.True(t, strings.Contains(bodyStr, `name="client_id" value="1111-2222-3333333-4444444"`))
	assert.True(t, strings.Contains(bodyStr, ` <h2>fight club`))
	assert.True(t, strings.Contains(bodyStr, `name="response_type" value="token"`))
	assert.True(t, strings.Contains(bodyStr, `<input type="hidden" name="scope" value="admin"`))

}

func TestHandleAuthorizeWithInvalidScope(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	var redirectCalled = false
	rs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		redirectCalled = true
		assert.True(t, strings.Contains(r.RequestURI, "error=access_denied&error_description=scope-problem"))
	}))
	defer rs.Close()

	returnVal := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		ClientID:        "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		ClientSecret:    "not for browser clients",
		RedirectURI:     rs.URL + "/foo",
		LoginProvider:   "xtrac://localhost:9000",
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	TestHTTPGet(t, addr+"/oauth2/authorize?client_id=1111-2222-3333333-4444444&redirect_uri="+rs.URL+"/foo&response_type=token&scope=invalid-scope", nil)
	appRepoMock.AssertCalled(t, "SystemRetrieveApplication", "1111-2222-3333333-4444444")

	assert.True(t, redirectCalled)

}

func TestHandleAuthorizeBadRedirectParam(t *testing.T) {
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

	resp := TestHTTPGet(t, addr+"/oauth2/authorize?client_id=1111-2222-3333333-4444444&redirect_uri=not-in-the-face&response_type=token", nil)
	appRepoMock.AssertCalled(t, "SystemRetrieveApplication", "1111-2222-3333333-4444444")
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAuthValidateBadClientId(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplication", "111-22-33").Return(nil, nil)

	resp, err := http.PostForm(addr+"/oauth2/validate",
		url.Values{"username": {"x"},
			"password":      {"y"},
			"authorize":     {"allow"},
			"response_type": {"token"},
			"client_id":     {"111-22-33"}})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "Invalid client id"))
}

func TestAuthValidateDenied(t *testing.T) {
	//TODO - use a second callback where we serve up a script to extract the page details sent
	//on deny and post those details to another test server.
	var callbackInvoked = false

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callbackInvoked = true
	}))
	defer ts.Close()

	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	returnVal := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		ClientID:        "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		ClientSecret:    "not for browser clients",
		RedirectURI:     ts.URL,
		LoginProvider:   "xtrac://localhost:9000",
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	_, err := http.PostForm(addr+"/oauth2/validate",
		url.Values{"username": {"x"},
			"password":      {"y"},
			"authorize":     {"deny"},
			"response_type": {"token"},
			"client_id":     {"1111-2222-3333333-4444444"}})
	assert.Nil(t, err)
	assert.True(t, callbackInvoked)
}

func TestAuthValidateAuthenticateFail(t *testing.T) {

	var loginCalled = false
	ls := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loginCalled = true
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ls.Close()

	//TODO - use a second callback where we serve up a script to extract the page details sent
	//on deny and post those details to another test server.
	var callbackInvoked = false

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callbackInvoked = true
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
		RedirectURI:     ts.URL,
		LoginProvider:   "xtrac://" + lsURL.Host,
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	_, err := http.PostForm(addr+"/oauth2/validate",
		url.Values{"username": {"x"},
			"password":      {"y"},
			"authorize":     {"allow"},
			"response_type": {"token"},
			"client_id":     {"1111-2222-3333333-4444444"}})
	assert.Nil(t, err)
	assert.True(t, callbackInvoked)
	assert.True(t, loginCalled)
}

func TestAuthValidateAuthenticateOkSecretsFail(t *testing.T) {

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
		RedirectURI:     ts.URL,
		LoginProvider:   "xtrac://" + lsURL.Host,
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	secretsMock := coreConfig.SecretsRepo.(*mocks.SecretsRepo)
	secretsMock.On("RetrievePrivateKeyForApp", "1111-2222-3333333-4444444").Return("", errors.New("Drat"))

	_, err := http.PostForm(addr+"/oauth2/validate",
		url.Values{"username": {"x"},
			"password":      {"y"},
			"authorize":     {"allow"},
			"response_type": {"token"},
			"client_id":     {"1111-2222-3333333-4444444"}})
	assert.Nil(t, err)
	assert.False(t, callbackInvoked)
	assert.True(t, loginCalled)
}

func TestAuthValidateAuthenticateOk(t *testing.T) {

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
		RedirectURI:     ts.URL,
		LoginProvider:   "xtrac://" + lsURL.Host,
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	privateKey, _, err := secrets.GenerateKeyPair()
	assert.Nil(t, err)

	secretsMock := coreConfig.SecretsRepo.(*mocks.SecretsRepo)
	secretsMock.On("RetrievePrivateKeyForApp", "1111-2222-3333333-4444444").Return(privateKey, nil)

	_, err = http.PostForm(addr+"/oauth2/validate",
		url.Values{"username": {"x"},
			"password":      {"y"},
			"authorize":     {"allow"},
			"response_type": {"token"},
			"client_id":     {"1111-2222-3333333-4444444"}})
	assert.Nil(t, err)
	assert.True(t, callbackInvoked)
	assert.True(t, loginCalled)
}

func TestImplGrantAuthValidateAuthenticateOkAdminScope(t *testing.T) {

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
		println(r.URL.RawQuery)
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
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	privateKey, pk, err := secrets.GenerateKeyPair()
	assert.Nil(t, err)

	secretsMock := coreConfig.SecretsRepo.(*mocks.SecretsRepo)
	secretsMock.On("RetrievePrivateKeyForApp", "1111-2222-3333333-4444444").Return(privateKey, nil)
	secretsMock.On("RetrievePublicKeyForApp", "1111-2222-3333333-4444444").Return(pk, nil)

	adminRepoMock := coreConfig.AdminRepo.(*mocks.AdminRepo)
	adminRepoMock.On("IsAdmin", "x").Return(true, nil)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			fmt.Println(req.URL.Fragment)
			m, _ := url.ParseQuery(req.URL.Fragment)

			accessToken := m.Get("access_token")
			token, err := jwt.Parse(accessToken, rolltoken.GenerateKeyExtractionFunction(core.SecretsRepo))
			assert.Nil(t, err)
			scope, ok := token.Claims["scope"].(string)
			assert.True(t, ok)
			assert.Equal(t, "admin", scope)

			assert.Equal(t, "x", token.Claims["sub"].(string))
			assert.Equal(t, "1111-2222-3333333-4444444", token.Claims["aud"].(string))

			assert.Equal(t, "Bearer", m.Get("token_type"))

			return nil
		},
	}

	_, err = client.PostForm(addr+"/oauth2/validate",
		url.Values{"username": {"x"},
			"password":      {"y"},
			"authorize":     {"allow"},
			"response_type": {"token"},
			"scope":         {"admin"},
			"client_id":     {"1111-2222-3333333-4444444"}})
	assert.Nil(t, err)
	assert.True(t, callbackInvoked)
	assert.True(t, loginCalled)
}
