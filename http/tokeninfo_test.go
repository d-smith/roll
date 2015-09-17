package http
import (
	"testing"
	"net/http"
	"github.com/xtracdev/xavi/Godeps/_workspace/src/github.com/stretchr/testify/assert"
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



