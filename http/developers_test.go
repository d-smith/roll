package http

import (
	"github.com/stretchr/testify/assert"
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
		Email:     "foo@gmail.com",
	}

	devRepoMock := coreConfig.DeveloperRepo.(*mocks.DeveloperRepo)
	devRepoMock.On("StoreDeveloper", &dev).Return(nil)

	resp := testHTTPPut(t, addr+"/v1/developers/foo@gmail.com", dev)
	devRepoMock.AssertCalled(t, "StoreDeveloper", &dev)

	checkResponseStatus(t, resp, http.StatusNoContent)
}

func TestGetDeveloper(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	devRepoMock := coreConfig.DeveloperRepo.(*mocks.DeveloperRepo)
	devRepoMock.On("RetrieveDeveloper", "joe@dev.com").Return(&roll.Developer{FirstName: "Joe", LastName: "Dev", Email: "joe@dev.com"}, nil)

	resp := testHTTPGet(t, addr+"/v1/developers/joe@dev.com", nil)
	devRepoMock.AssertCalled(t, "RetrieveDeveloper", "joe@dev.com")

	var actual roll.Developer

	checkResponseBody(t, resp, &actual)
	assert.Equal(t, "Joe", actual.FirstName)
	assert.Equal(t, "Dev", actual.LastName)
	assert.Equal(t, "joe@dev.com", actual.Email)

}
