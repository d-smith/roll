package http

import (
	"github.com/xtraclabs/roll/roll"
	"net/http"
	"testing"
)

func TestStoreDeveloper(t *testing.T) {
	core := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	dev := roll.Developer{
		FirstName: "Joe",
		LastName:  "Developer",
	}

	resp := testHttpPut(t, addr+"/v1/developers/foo@gmail.com", dev)

	checkResponseStatus(t, resp, http.StatusNoContent)
}
