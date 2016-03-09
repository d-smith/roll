package http

import (
	"bytes"
	"encoding/json"
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/roll/roll/mocks"
	"github.com/xtraclabs/rollsecrets/secrets"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

//Look at https://github.com/d-smith/go-examples/tree/master/certs and run the
//program to see how the public key was extracted from the certPEM below

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

//Look at and run https://github.com/d-smith/go-examples/tree/master/jwt/jwtkeycert to
//see where the assertion used in this test came from.
//Decoded assertion:
/*
{
  "aud": "captive",
  "iss": "1111-2222-3333333-4444444",
  "sub": "drscan"
}
*/
const jwtAssertion = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJjYXB0aXZlIiwiaXNzIjoiMTExMS0yMjIyLTMzMzMzMzMtNDQ0NDQ0NCIsInNjb3BlIjoiYSBiIGMgbm90YWRtaW4iLCJzdWIiOiJmb28ifQ.hVQ__RRprRNd09WK58GFZDw0pJEf5GebSPOs_fWXarZJEZN3UyTJKMNdjVOZVgwA5tJBp_hngsI8o4-ni5vhAPRNIXwAJhq2CuLGcy_F552FcW6xpSa3VM3rZgXXFmHq2_VurKrLwi0rX1I8Hax40mkwM797IOmzRyp89zfEQFpBybSJlYaGkcBdwpym2NSL7yrCBvWKZkidms3SgTOvh3XEF6FrXDWJrZuYGFjm0bStoYRvg_Zfy6JlkPOzHaCD5dIkdajqs3E5LP6bkBsSXK7m1228mFneQMxmi_CPhDqoIR3qpodptlt1gkLporBrVYbOk5dUBRmx3Of_D37knQ`

/*
{
  "aud": "captive"
  "iss": "1111-2222-3333333-4444444",
  "scope": "a b c admin",
  "sub": "foo"
}
*/
const jwtWithAdminScope = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJjYXB0aXZlIiwiaXNzIjoiMTExMS0yMjIyLTMzMzMzMzMtNDQ0NDQ0NCIsInNjb3BlIjoiYSBiIGMgYWRtaW4iLCJzdWIiOiJmb28ifQ.btsdQGbqk_5bOGlbBFXKVbrM1hInsH_EDnvXaOlpvq5Za96UZIRK1-gm5fkuAeCTW90rHKo1wxxMF4DF-nqlq9-TYy-YqP4Otnlrg3vDZ_8L1p4fiFHu1ktKPa4pOIVOWD74Vme2qa6xKlgGIzhVsZhZWxv-qFrwMjzSR9XZEJcw_KlYHQ89RgMNbeNL8gQ0kZpuWlmjBZkDtLM6hHJqL46HTCTLpcxzj-GfOMjzXT8MbzfxOJGKXhJH9lAhHKKYP5FwGiU0oEoQJ5OylnIrHinmwOypbXMIKl9ASYyp-0QLAuUcCIFzitf-coCwiESc3rq9Uka1om83X3KdVW-dmg`

func TestJWTFlowSetupMalformedPayload(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	buf := new(bytes.Buffer)
	buf.WriteString(`{"this won't parse`)

	resp := TestHTTPPutWithRollSubject(t, addr+"/v1/jwtflowcerts/11-22-33", buf)
	checkResponseStatus(t, resp, http.StatusBadRequest)
}

func TestJWTFlowSetupMissingClientSecret(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	requestBody := CertPutCtx{
		ClientSecret: "",
		CertPEM:      "yeah",
		CertIssuer:   "joe",
		CertAudience: "susie",
	}

	resp := TestHTTPPutWithRollSubject(t, addr+JWTFlowCertsURI+"123", requestBody)

	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "Request has empty ClientSecret"))
}

func TestJWTFlowSetupMissingCertPEM(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	requestBody := CertPutCtx{
		ClientSecret: "password123",
		CertPEM:      "",
		CertIssuer:   "joe",
		CertAudience: "susie",
	}

	resp := TestHTTPPutWithRollSubject(t, addr+JWTFlowCertsURI+"123", requestBody)

	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "Request has empty CertPEM"))
}

func TestJWTFlowSetupAppLookupError(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(nil, errors.New("Drat"))

	requestBody := CertPutCtx{
		ClientSecret: "foo",
		CertPEM:      "xxxxxx",
		CertIssuer:   "joe",
		CertAudience: "susie",
	}

	resp := TestHTTPPutWithRollSubject(t, addr+JWTFlowCertsURI+"1111-2222-3333333-4444444", requestBody)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

}

func TestJWTFlowSetupAppNotFound(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(nil, nil)

	requestBody := CertPutCtx{
		ClientSecret: "foo",
		CertPEM:      "xxxxxx",
		CertIssuer:   "joe",
		CertAudience: "susie",
	}

	resp := TestHTTPPutWithRollSubject(t, addr+JWTFlowCertsURI+"1111-2222-3333333-4444444", requestBody)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestJWTFlowSetupInvalidClientSecret(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	returnVal := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		ClientID:        "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		ClientSecret:    "not for browser clients",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	requestBody := CertPutCtx{
		ClientSecret: "foo",
		CertPEM:      "xxxxxx",
		CertIssuer:   "joe",
		CertAudience: "susie",
	}

	resp := TestHTTPPutWithRollSubject(t, addr+JWTFlowCertsURI+"1111-2222-3333333-4444444", requestBody)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestJWTFlowSetupInvalidCertPEM(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	returnVal := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		ClientID:        "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		ClientSecret:    "foo",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	requestBody := CertPutCtx{
		ClientSecret: "foo",
		CertPEM:      "xxxxxx",
		CertIssuer:   "joe",
		CertAudience: "susie",
	}

	resp := TestHTTPPutWithRollSubject(t, addr+JWTFlowCertsURI+"1111-2222-3333333-4444444", requestBody)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestJWTFlowSetupAppUpdateError(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	returnVal := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		ClientID:        "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		ClientSecret:    "foo",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	storeVal := roll.Application{
		DeveloperEmail:   "doug@dev.com",
		ClientID:         "1111-2222-3333333-4444444",
		ApplicationName:  "fight club",
		ClientSecret:     "foo",
		RedirectURI:      "http://localhost:3000/ab",
		LoginProvider:    "xtrac://localhost:9000",
		JWTFlowPublicKey: publicKey,
		JWTFlowIssuer:    "joe",
		JWTFlowAudience:  "susie",
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)
	appRepoMock.On("UpdateApplication", &storeVal, "rolltest").Return(errors.New("Ummm"))

	requestBody := CertPutCtx{
		ClientSecret: "foo",
		CertPEM:      certPEM,
		CertIssuer:   "joe",
		CertAudience: "susie",
	}

	resp := TestHTTPPutWithRollSubject(t, addr+JWTFlowCertsURI+"1111-2222-3333333-4444444", requestBody)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestJWTFlowSetupAppUpdateOk(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	returnVal := roll.Application{
		DeveloperEmail:  "doug@dev.com",
		ClientID:        "1111-2222-3333333-4444444",
		ApplicationName: "fight club",
		ClientSecret:    "foo",
		RedirectURI:     "http://localhost:3000/ab",
		LoginProvider:   "xtrac://localhost:9000",
	}

	storeVal := roll.Application{
		DeveloperEmail:   "doug@dev.com",
		ClientID:         "1111-2222-3333333-4444444",
		ApplicationName:  "fight club",
		ClientSecret:     "foo",
		RedirectURI:      "http://localhost:3000/ab",
		LoginProvider:    "xtrac://localhost:9000",
		JWTFlowPublicKey: publicKey,
		JWTFlowIssuer:    "joe",
		JWTFlowAudience:  "susie",
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)
	appRepoMock.On("UpdateApplication", &storeVal, "rolltest").Return(nil)

	requestBody := CertPutCtx{
		ClientSecret: "foo",
		CertPEM:      certPEM,
		CertIssuer:   "joe",
		CertAudience: "susie",
	}

	resp := TestHTTPPutWithRollSubject(t, addr+JWTFlowCertsURI+"1111-2222-3333333-4444444", requestBody)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestJWTFlowMissingAssertion(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"}})

	assert.Nil(t, err)
	body := responseAsString(t, resp)
	assert.True(t, strings.Contains(body, "assertion missing from request"))
}

func TestJWTFlowMalformedAssertion(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"},
			"assertion": {"this is not a jwt"}})

	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestJWTFlowValidAssertionAppLookupError(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplicationByJWTFlowAudience", "captive").Return(nil, errors.New("Drat"))

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"},
			"assertion": {jwtAssertion}})

	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestJWTFlowValidAssertionAppLookupReturnsNil(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplicationByJWTFlowAudience", "captive").Return(nil, nil)

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"},
			"assertion": {jwtAssertion}})

	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestJWTFlowValidAssertionOkYeah(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	returnVal := roll.Application{
		DeveloperEmail:   "doug@dev.com",
		ClientID:         "1111-2222-3333333-4444444",
		ApplicationName:  "fight club",
		ClientSecret:     "not for browser clients",
		RedirectURI:      "http://localhost:3000/ab",
		LoginProvider:    "xtrac://localhost:9000",
		JWTFlowPublicKey: publicKey,
		JWTFlowIssuer:    "1111-2222-3333333-4444444",
		JWTFlowAudience:  "captive",
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplicationByJWTFlowAudience", "captive").Return(&returnVal, nil)

	privateKey, publicKey, err := secrets.GenerateKeyPair()
	assert.Nil(t, err)

	secretsMock := coreConfig.SecretsRepo.(*mocks.SecretsRepo)
	secretsMock.On("RetrievePrivateKeyForApp", "1111-2222-3333333-4444444").Return(privateKey, nil)
	secretsMock.On("RetrievePublicKeyForApp", "1111-2222-3333333-4444444").Return(publicKey, nil)

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"},
			"assertion": {jwtAssertion}})

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bodyStr := responseAsString(t, resp)
	println(bodyStr)

	var jsonResponse accessTokenResponse
	err = json.Unmarshal([]byte(bodyStr), &jsonResponse)
	assert.Nil(t, err)
	assert.True(t, jsonResponse.AccessToken != "")
	assert.True(t, jsonResponse.TokenType == "Bearer")

	token, err := jwt.Parse(jsonResponse.AccessToken, roll.GenerateKeyExtractionFunction(core.SecretsRepo))
	assert.Nil(t, err)
	assert.Equal(t, "1111-2222-3333333-4444444", token.Claims["aud"].(string))
	assert.Equal(t, "foo", token.Claims["sub"].(string))
	assert.Equal(t, "", token.Claims["scope"].(string))
}

func TestJWTFlowValidAssertionOkAdminScope(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	returnVal := roll.Application{
		DeveloperEmail:   "doug@dev.com",
		ClientID:         "1111-2222-3333333-4444444",
		ApplicationName:  "fight club",
		ClientSecret:     "not for browser clients",
		RedirectURI:      "http://localhost:3000/ab",
		LoginProvider:    "xtrac://localhost:9000",
		JWTFlowPublicKey: publicKey,
		JWTFlowIssuer:    "1111-2222-3333333-4444444",
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplicationByJWTFlowAudience", "captive").Return(&returnVal, nil)

	privateKey, publicKey, err := secrets.GenerateKeyPair()
	assert.Nil(t, err)

	secretsMock := coreConfig.SecretsRepo.(*mocks.SecretsRepo)
	secretsMock.On("RetrievePrivateKeyForApp", "1111-2222-3333333-4444444").Return(privateKey, nil)
	secretsMock.On("RetrievePublicKeyForApp", "1111-2222-3333333-4444444").Return(publicKey, nil)

	resp, err := http.PostForm(addr+OAuth2TokenBaseURI,
		url.Values{"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"},
			"assertion": {jwtWithAdminScope}})

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bodyStr := responseAsString(t, resp)
	println(bodyStr)

	var jsonResponse accessTokenResponse
	err = json.Unmarshal([]byte(bodyStr), &jsonResponse)
	assert.Nil(t, err)
	assert.True(t, jsonResponse.AccessToken != "")
	assert.True(t, jsonResponse.TokenType == "Bearer")

	token, err := jwt.Parse(jsonResponse.AccessToken, roll.GenerateKeyExtractionFunction(core.SecretsRepo))
	assert.Nil(t, err)
	assert.Equal(t, "1111-2222-3333333-4444444", token.Claims["aud"].(string))
	assert.Equal(t, "foo", token.Claims["sub"].(string))
	assert.Equal(t, "admin", token.Claims["scope"].(string))
}

func TestJWTFlowGetResourceNotSpecified(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp := TestHTTPGetWithRollSubject(t, addr+JWTFlowCertsURI, nil)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestJWTFlowGetCertNotFound(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(nil, nil)

	resp := TestHTTPGetWithRollSubject(t, addr+JWTFlowCertsURI+"1111-2222-3333333-4444444", nil)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestJWTFlowGetCertRetrievalError(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(nil, errors.New("Drat"))

	resp := TestHTTPGetWithRollSubject(t, addr+JWTFlowCertsURI+"1111-2222-3333333-4444444", nil)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestJWTFlowGetCertOK(t *testing.T) {
	core, coreConfig := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	returnVal := roll.Application{
		DeveloperEmail:   "doug@dev.com",
		ClientID:         "1111-2222-3333333-4444444",
		ApplicationName:  "fight club",
		ClientSecret:     "not for browser clients",
		RedirectURI:      "http://localhost:3000/ab",
		LoginProvider:    "xtrac://localhost:9000",
		JWTFlowPublicKey: publicKey,
	}

	appRepoMock := coreConfig.ApplicationRepo.(*mocks.ApplicationRepo)
	appRepoMock.On("SystemRetrieveApplication", "1111-2222-3333333-4444444").Return(&returnVal, nil)

	resp := TestHTTPGetWithRollSubject(t, addr+JWTFlowCertsURI+"1111-2222-3333333-4444444", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var pkResp publicKeyCtx
	checkResponseBody(t, resp, &pkResp)
	assert.Equal(t, pkResp.PublicKey, publicKey)

}
