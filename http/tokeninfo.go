package http
import (
"github.com/xtraclabs/roll/roll"
"net/http"
"errors"
	"log"
jwt "github.com/dgrijalva/jwt-go"
	"fmt"
)


const (
	//TokenInfoURI is the base uri for the token validation service.
	TokenInfoURI = "/oauth2/tokeninfo"
)

func handleTokenInfo(core *roll.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleTokenInfoGet(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, errors.New("Method not allowed"))
		}
	})
}

func handleTokenInfoGet(core *roll.Core, w http.ResponseWriter, r *http.Request) {
	//Grab the token
	tokenString := r.FormValue("access_token")
	if tokenString == "" {
		log.Println("missing access token")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Parse and validate the token
	token, err := jwt.Parse(tokenString, roll.GenerateKeyExtractionFunction(core.SecretsRepo))
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Panic is there's no aud claim
	audience := token.Claims["aud"]
	if audience == "" {
		panic(errors.New("No aud claim in token"))
	}

	//Return the token info
	tokenInfo := fmt.Sprintf(`{"audience":"%s"}`, audience)
	w.Write([]byte(tokenInfo))
}