package http
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"net/http"
)

func TestRequiredQueryParamsPresent(t *testing.T) {
	t.Log("given a request with none of the required query params")
	t.Log("when requiredQueryParamsPresent is called")
	t.Log("then false is returned")

	req, _ := http.NewRequest("GET","/", nil)
	assert.False(t, requiredQueryParamsPresent(req))

	t.Log("given a request with some but not all of the required query params")
	t.Log("when requiredQueryParamsPresent is called")
	t.Log("then false is returned")
	req, _ = http.NewRequest("GET","/?client_id=123", nil)
	assert.False(t, requiredQueryParamsPresent(req))

	t.Log("given a request with all of the required query params")
	t.Log("when requiredQueryParamsPresent is called")
	t.Log("then true is returned")
	req, _ = http.NewRequest("POST","/?client_id=123&redirect_uri=x&response_type=X", nil)
	assert.True(t, requiredQueryParamsPresent(req))
}
