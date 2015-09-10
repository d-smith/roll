package http
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"net/http"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/roll/roll/mocks"
	"errors"
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
