package http

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/roll/roll/mocks"
	"net/http"
	"testing"
)

func TestStoreDeveloperOK(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	dev := roll.Developer{
		FirstName: "Joe",
		LastName:  "Developer",
		Email:     "foo@gmail.com",
		ID:        "rolltest",
	}

	devRepoMock := coreConfig.DeveloperRepo.(*mocks.DeveloperRepo)
	devRepoMock.On("StoreDeveloper", &dev).Return(nil)

	resp := TestHTTPPutWithRollSubject(t, addr+"/v1/developers/foo@gmail.com", dev)
	devRepoMock.AssertCalled(t, "StoreDeveloper", &dev)

	checkResponseStatus(t, resp, http.StatusNoContent)
}

func TestStoreDeveloperStorageFault(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	dev := roll.Developer{
		FirstName: "Joe",
		LastName:  "Developer",
		Email:     "foo@gmail.com",
		ID:        "rolltest",
	}

	devRepoMock := coreConfig.DeveloperRepo.(*mocks.DeveloperRepo)
	devRepoMock.On("StoreDeveloper", &dev).Return(errors.New("can't store"))

	resp := TestHTTPPutWithRollSubject(t, addr+"/v1/developers/foo@gmail.com", dev)
	devRepoMock.AssertCalled(t, "StoreDeveloper", &dev)

	checkResponseStatus(t, resp, http.StatusInternalServerError)
}

func TestStoreDeveloperInvalidEmailResource(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp := TestHTTPPutWithRollSubject(t, addr+"/v1/developers/<script/>", nil)
	checkResponseStatus(t, resp, http.StatusBadRequest)
}

func TestStoreDevBodyParseError(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	buf := new(bytes.Buffer)
	buf.WriteString(`{"this won't parse`)

	resp := TestHTTPPutWithRollSubject(t, addr+"/v1/developers/foo@dev.com", buf)

	checkResponseStatus(t, resp, http.StatusBadRequest)

}

func TestStoreDeveloperInvalidContent(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	dev := roll.Developer{
		FirstName: "Joe<script>",
		LastName:  "Developer",
		Email:     "foo@gmail.com",
	}

	resp := TestHTTPPutWithRollSubject(t, addr+"/v1/developers/foo@gmail.com", dev)

	checkResponseStatus(t, resp, http.StatusBadRequest)
}

func TestGetDeveloperInvalidEmailResource(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp := TestHTTPGetWithRollSubject(t, addr+"/v1/developers/<script/>", nil)
	checkResponseStatus(t, resp, http.StatusBadRequest)
}

func TestDeveloperUnsupportedMethod(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp := TestHTTPPostWithRollSubject(t, addr+"/v1/developers/1111-2222-3333333-4444444", nil)
	checkResponseStatus(t, resp, http.StatusMethodNotAllowed)
}

func TestGetDeveloperOk(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	devRepoMock := coreConfig.DeveloperRepo.(*mocks.DeveloperRepo)
	devRepoMock.On("RetrieveDeveloper", "joe@dev.com", "rolltest", false).Return(&roll.Developer{FirstName: "Joe", LastName: "Dev", Email: "joe@dev.com"}, nil)

	resp := TestHTTPGetWithRollSubject(t, addr+"/v1/developers/joe@dev.com", nil)
	devRepoMock.AssertCalled(t, "RetrieveDeveloper", "joe@dev.com", "rolltest", false)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var actual roll.Developer

	checkResponseBody(t, resp, &actual)
	assert.Equal(t, "Joe", actual.FirstName)
	assert.Equal(t, "Dev", actual.LastName)
	assert.Equal(t, "joe@dev.com", actual.Email)

}

func TestGetDeveloperMissingRequestContext(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp := TestHTTPGet(t, addr+"/v1/developers/joe@dev.com", nil)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestGetDevelopersListMissingRequestContext(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp := TestHTTPGet(t, addr+"/v1/developers", nil)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestGetDevelopersOk(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	devRepoMock := coreConfig.DeveloperRepo.(*mocks.DeveloperRepo)
	devs := []roll.Developer{
		roll.Developer{FirstName: "Joe", LastName: "Dev", Email: "joe@dev.com"},
		roll.Developer{FirstName: "Jill", LastName: "Dev", Email: "jill@dev.com"},
	}
	devRepoMock.On("ListDevelopers", "rolltest", false).Return(devs, nil)

	resp := TestHTTPGetWithRollSubject(t, addr+"/v1/developers", nil)
	devRepoMock.AssertCalled(t, "ListDevelopers", "rolltest", false)

	var actual []roll.Developer
	checkResponseBody(t, resp, &actual)
	assert.Equal(t, 2, len(actual))

	for _, d := range actual {
		switch d.Email {
		case "joe@dev.com":
			assert.Equal(t, "Joe", d.FirstName)
			assert.Equal(t, "Dev", d.LastName)
		case "jill@dev.com":
			assert.Equal(t, "Jill", d.FirstName)
			assert.Equal(t, "Dev", d.LastName)
		default:
			assert.Error(t, errors.New("Unexpected dev email"))
		}
	}

}

func TestGetDevelopersDBError(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	devRepoMock := coreConfig.DeveloperRepo.(*mocks.DeveloperRepo)
	devRepoMock.On("ListDevelopers", "rolltest", false).Return(nil, errors.New("db error"))

	resp := TestHTTPGetWithRollSubject(t, addr+"/v1/developers", nil)
	devRepoMock.AssertCalled(t, "ListDevelopers", "rolltest", false)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestGetDeveloperRetrieveError(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	devRepoMock := coreConfig.DeveloperRepo.(*mocks.DeveloperRepo)
	devRepoMock.On("RetrieveDeveloper", "joe@dev.com", "rolltest", false).Return(nil, errors.New("retrieve error"))

	resp := TestHTTPGetWithRollSubject(t, addr+"/v1/developers/joe@dev.com", nil)
	devRepoMock.AssertCalled(t, "RetrieveDeveloper", "joe@dev.com", "rolltest", false)

	checkResponseStatus(t, resp, http.StatusInternalServerError)
}

func TestGetNonExistentDeveloper(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	devRepoMock := coreConfig.DeveloperRepo.(*mocks.DeveloperRepo)
	devRepoMock.On("RetrieveDeveloper", "joe@dev.com", "rolltest", false).Return(nil, nil)

	resp := TestHTTPGetWithRollSubject(t, addr+"/v1/developers/joe@dev.com", nil)
	devRepoMock.AssertCalled(t, "RetrieveDeveloper", "joe@dev.com", "rolltest", false)

	checkResponseStatus(t, resp, http.StatusNotFound)

}

func TestCheckResponseStatus(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	devRepoMock := coreConfig.DeveloperRepo.(*mocks.DeveloperRepo)
	devRepoMock.On("RetrieveDeveloper", "joe@dev.com", "rolltest", false).Return(nil, nil)

	resp := TestHTTPGetWithRollSubject(t, addr+"/v1/developers/joe@dev.com", nil)
	devRepoMock.AssertCalled(t, "RetrieveDeveloper", "joe@dev.com", "rolltest", false)

	ok := checkResponseStatus(nil, resp, http.StatusOK)
	assert.False(t, ok)
}
