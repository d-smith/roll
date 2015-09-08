package http
import (
	"github.com/xtraclabs/roll/roll"
	"net/http"
	"errors"
)

const (
	OAuth2TokenBaseURI = "/oauth2/token"
)

func handleToken(core *roll.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			handleTokenPost(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, errors.New("Method not allowed"))
		}
	})
}

func handleTokenPost(core *roll.Core, w http.ResponseWriter, r *http.Request) {
	//Verify the form params are as expected: grant_type is authorization_code,
	//a code is present, client_id and client_secret are provided, redirect_uri is
	//provided. The content type should be application/x-www-form-urlencoded

	//Verify the client id and secret, plus the redirect_uri by doing a lookup
	//of the app by client id as api key value

	//Validate the code - it should be a token signed with the users' private key

	//If everything is cool, generate a JWT access token

	//Respond with a JSON document included the access_token and a token type of
	//bearer
	respondError(w, http.StatusInternalServerError, errors.New("Not implemented"))

}
