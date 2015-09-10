package http
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"net/http"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/roll/roll/mocks"
	"errors"
	"net/http/httptest"
	"strings"
	"io/ioutil"
)

func TestRequiredQueryParamsPresent(t *testing.T) {
	t.Log("given a request with none of the required query params")
	t.Log("when requiredQueryParamsPresent is called")
	t.Log("then false is returned")

	req, _ := http.NewRequest("GET","/", nil)
	assert.False(t, requiredQueryParamsPresent(req))

	t.Log("given a request with some but not all of the required query params")
	t.Log("when requiredQueryParamsPresent is called")
	t.Log("then false is returned")
	req, _ = http.NewRequest("GET","/?client_id=123", nil)
	assert.False(t, requiredQueryParamsPresent(req))

	t.Log("given a request with all of the required query params")
	t.Log("when requiredQueryParamsPresent is called")
	t.Log("then true is returned")
	req, _ = http.NewRequest("POST","/?client_id=123&redirect_uri=x&response_type=X", nil)
	assert.True(t, requiredQueryParamsPresent(req))
}

func TestInputParamsValid(t *testing.T) {
	core, coreConfig := NewTestCore()

	returnVal := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		APIKey:          "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		APISecret:       "not for browser clients",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)


	req, _ := http.NewRequest("POST","/?client_id=1111-2222-3333333-4444444&redirect_uri=http://localhost:3000/ab&response_type=token", nil)

	app, err := validateInputParams(core, req)
	assert.Nil(t, err)
	assert.NotNil(t, app)
}

func TestInputParamsInvalidResponseType(t *testing.T) {
	req, _ := http.NewRequest("POST","/?client_id=1111-2222-3333333-4444444&redirect_uri=http://localhost:3000/ab&response_type=bad", nil)
	app, err := validateInputParams(nil, req)
	assert.Nil(t,app)
	assert.NotNil(t,err)
}

func TestInputParamsMissingResponseType(t *testing.T) {
	req, _ := http.NewRequest("POST","/?client_id=1111-2222-3333333-4444444&redirect_uri=http://localhost:3000/ab", nil)
	app, err := validateInputParams(nil, req)
	assert.Nil(t,app)
	assert.NotNil(t,err)
}

func TestInputParamsMissingClientID(t *testing.T) {
	core, coreConfig := NewTestCore()
	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "").Return(nil, nil)

	req, _ := http.NewRequest("POST","/?redirect_uri=http://localhost:3000/ab&response_type=code", nil)
	app, err := validateInputParams(core, req)
	assert.Nil(t,app)
	assert.NotNil(t,err)
	println(err.Error())
}

func TestInputParamsNoSuchClientId(t *testing.T) {
	core, coreConfig := NewTestCore()

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(nil, errors.New("whoops"))

	req, _ := http.NewRequest("POST","/?client_id=1111-2222-3333333-4444444&redirect_uri=http://localhost:3000/ab&response_type=token", nil)

	app, err := validateInputParams(core, req)
	assert.NotNil(t, err)
	assert.Nil(t, app)
}

func TestInputParamsInvalidRedirectURI(t *testing.T) {
	core, coreConfig := NewTestCore()

	returnVal := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		APIKey:          "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		APISecret:       "not for browser clients",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)


	req, _ := http.NewRequest("POST","/?client_id=1111-2222-3333333-4444444&redirect_uri=bogus&response_type=token", nil)
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

	req, _ := http.NewRequest("POST","/?client_id=1111-2222-3333333-4444444&redirect_uri=bogus&response_type=code", nil)

	err := executeAuthTemplate(w,req,pageCtx)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	body :=  w.Body.String()
	assert.True(t, strings.Contains(body, `name="client_id" value="test-application-key"`))
	assert.True(t, strings.Contains(body, ` <h2>test-application-name`))
	assert.True(t, strings.Contains(body, `name="response_type" value="code"`))
}

func TestExecuteAuthTemplateForToken(t *testing.T) {
	w := httptest.NewRecorder()
	pageCtx := &authPageContext{
		AppName:  "test-application-name",
		ClientID: "test-application-key",
	}

	req, _ := http.NewRequest("POST","/?client_id=1111-2222-3333333-4444444&redirect_uri=bogus&response_type=token", nil)

	err := executeAuthTemplate(w,req,pageCtx)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	body :=  w.Body.String()
	assert.True(t, strings.Contains(body, `name="client_id" value="test-application-key"`))
	assert.True(t, strings.Contains(body, ` <h2>test-application-name`))
	assert.True(t, strings.Contains(body, `name="response_type" value="token"`))
}

func TestExecuteAuthTemplateMissingResponseType(t *testing.T) {
	w := httptest.NewRecorder()
	pageCtx := &authPageContext{
		AppName:  "test-application-name",
		ClientID: "test-application-key",
	}

	req, _ := http.NewRequest("POST","/?client_id=1111-2222-3333333-4444444&redirect_uri=bogus", nil)

	err := executeAuthTemplate(w,req,pageCtx)
	assert.NotNil(t, err)

}

func TestExecuteAuthTemplateBogusResponseType(t *testing.T) {
	w := httptest.NewRecorder()
	pageCtx := &authPageContext{
		AppName:  "test-application-name",
		ClientID: "test-application-key",
	}

	req, _ := http.NewRequest("POST","/?client_id=1111-2222-3333333-4444444&redirect_uri=bogus&response_type=bogus", nil)

	err := executeAuthTemplate(w,req,pageCtx)
	assert.NotNil(t, err)

}

func TestHandleAuthorize(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	returnVal := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		APIKey:          "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		APISecret:       "not for browser clients",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)



	resp := testHTTPGet(t, addr+"/oauth2/authorize?client_id=1111-2222-3333333-4444444&redirect_uri=http://localhost:3000/ab&response_type=token", nil)
	appRepoMock.AssertCalled(t, "RetrieveApplication", "1111-2222-3333333-4444444")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bodyStr := string(body)
	assert.True(t, strings.Contains(bodyStr, `name="client_id" value="1111-2222-3333333-4444444"`))
	assert.True(t, strings.Contains(bodyStr, ` <h2>fight club`))
	assert.True(t, strings.Contains(bodyStr, `name="response_type" value="token"`))

}

func TestHandleAuthorizeUnsupportedMethod(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.Post(addr+"/oauth2/authorize","",nil)
	assert.Nil(t,err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestHandleAuthorizeMissingParams(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp := testHTTPGet(t, addr+"/oauth2/authorize",nil)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestHandleAuthorizeBadRedirectParam(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	returnVal := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		APIKey:          "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		APISecret:       "not for browser clients",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)



	resp := testHTTPGet(t, addr+"/oauth2/authorize?client_id=1111-2222-3333333-4444444&redirect_uri=not-in-the-face&response_type=token", nil)
	appRepoMock.AssertCalled(t, "RetrieveApplication", "1111-2222-3333333-4444444")
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
