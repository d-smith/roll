package http

import (
	"github.com/CenturyLinkLabs/imagelayers/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/roll/roll/mocks"
	"net/http"
	"testing"
)

func TestStoreDeveloper(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	dev := roll.Developer{
		FirstName: "Joe",
		LastName:  "Developer",
	}

	testObj := coreConfig.DeveloperRepo.(*mocks.DeveloperRepo)
	testObj.On("StoreDeveloper", &dev).Return(nil)

	resp := testHttpPut(t, addr+"/v1/developers/foo@gmail.com", dev)
	testObj.AssertCalled(t, "StoreDeveloper", &dev)

	checkResponseStatus(t, resp, http.StatusNoContent)
}

func TestGetDeveloper(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	testObj := coreConfig.DeveloperRepo.(*mocks.DeveloperRepo)
	testObj.On("RetrieveDeveloper", "joe@dev.com").Return(&roll.Developer{FirstName: "Joe", LastName: "Dev", Email: "joe@dev.com"}, nil)

	resp := testHttpGet(t, addr+"/v1/developers/joe@dev.com", nil)
	testObj.AssertCalled(t, "RetrieveDeveloper", "joe@dev.com")

	var actual roll.Developer

	checkResponseBody(t, resp, &actual)
	assert.Equal(t, "Joe", actual.FirstName)
	assert.Equal(t, "Dev", actual.LastName)
	assert.Equal(t, "joe@dev.com", actual.Email)

}
