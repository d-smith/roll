package http
import (
	"testing"
	"net/http"
	"github.com/xtracdev/xavi/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/roll/roll/mocks"
	"github.com/xtraclabs/roll/secrets"
	"encoding/json"
)

func TestMissingAccessToken(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.Get(addr + TokenInfoURI)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}


func TestInvalidAccessToken(t *testing.T) {
	core, _ := NewTestCore()
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.Get(addr + TokenInfoURI + "?access_token=xxx.xxx.xxx")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestValidAccessToken(t *testing.T) {
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

	privateKey, publicKey, err := secrets.GenerateKeyPair()
	assert.Nil(t, err)

	secretsMock := coreConfig.SecretsRepo.(*mocks.SecretsRepo)
	secretsMock.On("RetrievePrivateKeyForApp", "1111-2222-3333333-4444444").Return(privateKey, nil)
	secretsMock.On("RetrievePublicKeyForApp", "1111-2222-3333333-4444444").Return(publicKey, nil)

	token, err := roll.GenerateToken(&returnVal, privateKey)
	assert.Nil(t, err)

	resp, err := http.Get(addr + TokenInfoURI + "?access_token=" + token)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body := responseAsString(t, resp)
	var ti TokenInfo

	err = json.Unmarshal([]byte(body), &ti)
	assert.Nil(t,err)
	assert.Equal(t, "1111-2222-3333333-4444444", ti.Audience)

}

