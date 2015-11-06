package http

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/roll/roll/mocks"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"testing"
)

//These functions inspired by/adopted from github.com/hashicorp/vault testing and http_test files in the vault
//http package

//TestListener creates a lister on the localhost, selecting a free port to listen on.
func TestListener(t *testing.T) (net.Listener, string) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fail()
	}

	addr := "http://" + ln.Addr().String()
	return ln, addr
}

//TestServerWithListener creates a test server using a given listener, running it in a new go routine.
func TestServerWithListener(t *testing.T, ln net.Listener, core *roll.Core) {
	mux := http.NewServeMux()
	mux.Handle("/", Handler(core))

	server := &http.Server{
		Addr:    ln.Addr().String(),
		Handler: mux,
	}

	go server.Serve(ln)
}

//TestServer return a test listener and its address after firing up a
//TestServerWithListener
func TestServer(t *testing.T, core *roll.Core) (net.Listener, string) {
	ln, addr := TestListener(t)
	TestServerWithListener(t, ln, core)
	return ln, addr
}

type TestIDGen struct{}

func (tig TestIDGen) GenerateID() (string, error) {
	return "steve", nil
}

//NewTestCore returns a roll.Core instance with mocked implementations of its internal dependencies
func NewTestCore() (*roll.Core, *roll.CoreConfig) {
	var coreConfig = roll.CoreConfig{}
	coreConfig.DeveloperRepo = new(mocks.DeveloperRepo)
	coreConfig.ApplicationRepo = new(mocks.ApplicationRepo)
	coreConfig.AdminRepo = new(mocks.AdminRepo)
	coreConfig.SecretsRepo = new(mocks.SecretsRepo)
	coreConfig.IdGenerator = TestIDGen{}
	coreConfig.Secure = false
	return roll.NewCore(&coreConfig), &coreConfig
}

func checkFatal(t assert.TestingT, err error) {
	if err != nil {
		log.Println("checkFatal passed an error")
		assert.Fail(t, err.Error())
	}
}

func TestHTTPGet(t assert.TestingT, addr string, body interface{}) *http.Response {
	return testHTTPData(t, "GET", addr, false, body)
}

func TestHTTPGetWithRollSubject(t assert.TestingT, addr string, body interface{}) *http.Response {
	return testHTTPData(t, "GET", addr, true, body)
}

func TestHTTPPut(t assert.TestingT, addr string, body interface{}) *http.Response {
	return testHTTPData(t, "PUT", addr, false, body)
}

func TestHTTPPutWithRollSubject(t assert.TestingT, addr string, body interface{}) *http.Response {
	return testHTTPData(t, "PUT", addr, true, body)
}

func TestHTTPPost(t assert.TestingT, addr string, body interface{}) *http.Response {
	return testHTTPData(t, "POST", addr, false, body)
}

func TestHTTPPostWithRollSubject(t assert.TestingT, addr string, body interface{}) *http.Response {
	return testHTTPData(t, "POST", addr, true, body)
}

func testHTTPData(t assert.TestingT, method string, addr string, rollSubject bool, body interface{}) *http.Response {
	bodyReader := new(bytes.Buffer)
	if body != nil {
		enc := json.NewEncoder(bodyReader)
		err := enc.Encode(body)
		checkFatal(t, err)
	}

	req, err := http.NewRequest(method, addr, bodyReader)
	checkFatal(t, err)

	req.Header.Set("Content-Type", "application/json")
	if rollSubject {
		req.Header.Set("X-Roll-Subject", "rolltest")
	}

	client := http.DefaultClient

	resp, err := client.Do(req)
	checkFatal(t, err)

	return resp
}

func checkResponseStatus(t *testing.T, resp *http.Response, code int) bool {
	var ok = true
	if resp.StatusCode != code {
		ok = false
		body := new(bytes.Buffer)
		io.Copy(body, resp.Body)
		resp.Body.Close()

		if t != nil {
			t.Errorf("Expected status %d got %d with body \n%s\n", code, resp.StatusCode, body)
		} else {
			log.Printf("Expected status %d got %d with body \n%s\n", code, resp.StatusCode, body)
		}
	}

	return ok
}

func checkResponseBody(t *testing.T, resp *http.Response, out interface{}) {
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(out); err != nil {
		t.Fatalf("Error: %s", err)
	}
}

func responseAsString(t *testing.T, r *http.Response) string {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if assert.Nil(t, err) {
		return string(body)
	}

	return ""
}
