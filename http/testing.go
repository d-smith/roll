package http

import (
	"bytes"
	"encoding/json"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/roll/roll/mocks"
	"io"
	"net"
	"net/http"
	"testing"
)

//These functions inspired by/adopted from github.com/hashicorp/vault testing and http_test files in the vault
//http package

func TestListener(t *testing.T) (net.Listener, string) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fail()
	}

	addr := "http://" + ln.Addr().String()
	return ln, addr
}

func TestServerWithListener(t *testing.T, ln net.Listener, core *roll.Core) {
	mux := http.NewServeMux()
	mux.Handle("/", Handler(core))

	server := &http.Server{
		Addr:    ln.Addr().String(),
		Handler: mux,
	}

	go server.Serve(ln)
}

func TestServer(t *testing.T, core *roll.Core) (net.Listener, string) {
	ln, addr := TestListener(t)
	TestServerWithListener(t, ln, core)
	return ln, addr
}

func NewTestCore() (*roll.Core, *roll.CoreConfig) {
	var coreConfig = roll.CoreConfig{}
	coreConfig.DeveloperRepo = new(mocks.DeveloperRepo)
	return roll.NewCore(&coreConfig), &coreConfig
}

func checkFatal(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("error: %s", err)
	}
}

func testHttpGet(t *testing.T, addr string, body interface{}) *http.Response {
	return testHttpData(t, "GET", addr, body)
}

func testHttpPut(t *testing.T, addr string, body interface{}) *http.Response {
	return testHttpData(t, "PUT", addr, body)
}

func testHttpData(t *testing.T, method string, addr string, body interface{}) *http.Response {
	bodyReader := new(bytes.Buffer)
	if body != nil {
		enc := json.NewEncoder(bodyReader)
		err := enc.Encode(body)
		checkFatal(t, err)
	}

	req, err := http.NewRequest(method, addr, bodyReader)
	checkFatal(t, err)

	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient

	resp, err := client.Do(req)
	checkFatal(t, err)

	return resp
}

func checkResponseStatus(t *testing.T, resp *http.Response, code int) {
	if resp.StatusCode != code {
		body := new(bytes.Buffer)
		io.Copy(body, resp.Body)
		resp.Body.Close()

		t.Fatalf("Expected status %d got %d with body \n%s\n", code, resp.StatusCode, body)
	}
}

func checkResponseBody(t *testing.T, resp *http.Response, out interface{}) {
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(out); err != nil {
		t.Fatalf("Error: %s", err)
	}
}
