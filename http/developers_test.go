package http

import (
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
