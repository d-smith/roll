package http

import (
	"testing"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/roll/roll/mocks"
	"net/http"
	"github.com/stretchr/testify/assert"
)

func TestStoreApp(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	app := roll.Application{
		ApplicationName:"ambivilant birds",
		DeveloperEmail:"doug@dev.com",
		APIKey:"1111-2222-3333333-4444444",
	}

	testObj := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	testObj.On("StoreApplication", &app).Return(nil)

	resp := testHttpPut(t, addr+"/v1/applications/1111-2222-3333333-4444444", app)
	testObj.AssertCalled(t, "StoreApplication", &app)

	checkResponseStatus(t, resp, http.StatusNoContent)
}

func TestGetApplication(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	returnVal := roll.Application {
		DeveloperEmail:"doug@dev.com",
		APIKey: "1111-2222-3333333-4444444",
		ApplicationName:"fight club",
		APISecret: "not for browser clients",
	}

	testObj := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	testObj.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	resp := testHttpGet(t, addr+"/v1/applications/1111-2222-3333333-4444444", nil)
	testObj.AssertCalled(t, "RetrieveApplication", "1111-2222-3333333-4444444")

	var actual roll.Application

	checkResponseBody(t, resp, &actual)
	assert.Equal(t, "doug@dev.com", actual.DeveloperEmail)
	assert.Equal(t, "1111-2222-3333333-4444444", actual.APIKey)
	assert.Equal(t, "fight club", actual.ApplicationName)
	assert.Equal(t, "not for browser clients", actual.APISecret)

}
