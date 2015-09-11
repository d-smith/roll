package http

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/roll/roll/mocks"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

const certPEM = `
-----BEGIN CERTIFICATE-----
MIIC+DCCAeKgAwIBAgIRAIJaB8pAErenO9pMBUDo3awwCwYJKoZIhvcNAQELMBIx
EDAOBgNVBAoTB0FjbWUgQ28wHhcNMTUwODI4MTM0MzU3WhcNMTYwODI3MTM0MzU3
WjASMRAwDgYDVQQKEwdBY21lIENvMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB
CgKCAQEAwoH4nc/B3/i1D1TjCx3kgC6ygX3WHDv/xHtAoRAgHFUVElo3PznbxLAk
MvElVdAevCJaiuJaLiZARKLvwSJh08/9y+WMYa1nDjINk6UqG3huPXdJmTguzleO
c7UrCW4WKSo2HbeqYlF4BOiqnQhdDncUh5BgR8JXuiueMn2Ka59lkB/i+ryOt5W7
kaFKJhQEV67+fuES/5WfE+B4XsfT/ctXnGY0zrEInbJlyKwAzyCWJOJFrZte8cxs
235q3VMAhMRDU1IGNuWBIntfEXZgUXqI1Z9gsdbfTsQQ+xWhQCCOJwDrxAEg1Udk
dWn6NGWevsH4JoM9JzzOeSH8ZYPrVQIDAQABo00wSzAOBgNVHQ8BAf8EBAMCAKAw
EwYDVR0lBAwwCgYIKwYBBQUHAwEwDAYDVR0TAQH/BAIwADAWBgNVHREEDzANggtN
QUNMQjAxNTgwMzALBgkqhkiG9w0BAQsDggEBALJtJGaXx9At98CvEWKBpiGYqjUu
aiQHS5R61R/g8iqWkct77cqN6SBWTf138NZ3j3mvfROCoU96BEMEl0Fk9apLrikI
9Ns9/sl4nL1IOR56vddm46DfEV5CpMCAgrMGhFMJiaW4t9HvYjpBSs8T5n4tGqu/
JsvPhLGOcu5i4RiPpwM8f4fhnD3sija334jj5meJwg0NR8eO3ro1zaH+0hMQ7l8Q
tFJusSJenG28q9MXpOoCG6KLCmSCrIfDRYIpJQ0d5fXLO4YG92KFFqrf2ycOTydY
hN9G5ZWaErEY5j+sbYmeJBtEM5v6BQJotJh2SAh8RpYr69qJPLw6fdTu+mU=
-----END CERTIFICATE-----`

const publicKey = `-----BEGIN RSA PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwoH4nc/B3/i1D1TjCx3k
gC6ygX3WHDv/xHtAoRAgHFUVElo3PznbxLAkMvElVdAevCJaiuJaLiZARKLvwSJh
08/9y+WMYa1nDjINk6UqG3huPXdJmTguzleOc7UrCW4WKSo2HbeqYlF4BOiqnQhd
DncUh5BgR8JXuiueMn2Ka59lkB/i+ryOt5W7kaFKJhQEV67+fuES/5WfE+B4XsfT
/ctXnGY0zrEInbJlyKwAzyCWJOJFrZte8cxs235q3VMAhMRDU1IGNuWBIntfEXZg
UXqI1Z9gsdbfTsQQ+xWhQCCOJwDrxAEg1UdkdWn6NGWevsH4JoM9JzzOeSH8ZYPr
VQIDAQAB
-----END RSA PUBLIC KEY-----
`

func TestJWTFlowMissingClientSecret(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.PostForm(addr+JWTFlowCertsURI,
		url.Values{})

	assert.Nil(t, err)
	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "client_secret missing from request"))
}

func TestJWTFlowMissingCertPEM(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.PostForm(addr+JWTFlowCertsURI,
		url.Values{"client_secret": {"foo"}})

	assert.Nil(t, err)
	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "cert_pem missing from request"))
}

func TestJWTFlowAppLookupError(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(nil, errors.New("Drat"))

	resp, err := http.PostForm(addr+JWTFlowCertsURI+"1111-2222-3333333-4444444",
		url.Values{"client_secret": {"foo"},
			"cert_pem": {"xxxxxx"}})

	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

}

func TestJWTFlowAppNotFound(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(nil, nil)

	resp, err := http.PostForm(addr+JWTFlowCertsURI+"1111-2222-3333333-4444444",
		url.Values{"client_secret": {"foo"},
			"cert_pem": {"xxxxxx"}})

	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestJWTFlowInvalidClientSecret(t *testing.T) {
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

	resp, err := http.PostForm(addr+JWTFlowCertsURI+"1111-2222-3333333-4444444",
		url.Values{"client_secret": {"foo"},
			"cert_pem": {"xxxxxx"}})

	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestJWTFlowInvalidCertPEM(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	returnVal := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		APIKey:          "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		APISecret:       "foo",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	resp, err := http.PostForm(addr+JWTFlowCertsURI+"1111-2222-3333333-4444444",
		url.Values{"client_secret": {"foo"},
			"cert_pem": {"xxxxxx"}})

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestJWTFlowAppUpdateError(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	returnVal := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		APIKey:          "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		APISecret:       "foo",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	storeVal := roll.Application{
		DeveloperEmail:   "doug@dev.com",
		APIKey:           "1111-2222-3333333-4444444",
		ApplicationName:  "fight club",
		APISecret:        "foo",
		RedirectURI:      "http://localhost:3000/ab",
		LoginProvider:    "xtrac://localhost:9000",
		JWTFlowPublicKey: publicKey,
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)
	appRepoMock.On("StoreApplication", &storeVal).Return(errors.New("Ummm"))

	resp, err := http.PostForm(addr+JWTFlowCertsURI+"1111-2222-3333333-4444444",
		url.Values{"client_secret": {"foo"},
			"cert_pem": {certPEM}})

	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestJWTFlowAppUpdateOk(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	returnVal := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		APIKey:          "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		APISecret:       "foo",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	storeVal := roll.Application{
		DeveloperEmail:   "doug@dev.com",
		APIKey:           "1111-2222-3333333-4444444",
		ApplicationName:  "fight club",
		APISecret:        "foo",
		RedirectURI:      "http://localhost:3000/ab",
		LoginProvider:    "xtrac://localhost:9000",
		JWTFlowPublicKey: publicKey,
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("RetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)
	appRepoMock.On("StoreApplication", &storeVal).Return(nil)

	resp, err := http.PostForm(addr+JWTFlowCertsURI+"1111-2222-3333333-4444444",
		url.Values{"client_secret": {"foo"},
			"cert_pem": {certPEM}})

	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
