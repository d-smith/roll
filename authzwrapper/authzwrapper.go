package authzwrapper

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/xtraclabs/roll/roll"
	"log"
	"net/http"
	"strings"
)

type authHandler struct {
	handler     http.Handler
	secretsRepo roll.SecretsRepo
	whiteList   map[string]string
}

//Wrap takes a handler and decorates it with JWT bearer token validation.
func Wrap(secretsRepo roll.SecretsRepo, whitelistedClientIDs []string, h http.Handler) http.Handler {
	wl := make(map[string]string)
	for _, cid := range whitelistedClientIDs {
		wl[cid] = cid
	}

	return &authHandler{
		handler:     h,
		secretsRepo: secretsRepo,
		whiteList:   wl,
	}
}

//AddWhiteListedClientID adds a client id to the whitelist. When the whitelist contains 1 or more
//client IDs, the aud claim of the bearer token is checked against the whitelist -- if present
//in the list access is granted.
func (ah authHandler) AddWhiteListedClientID(clientId string) {
	ah.whiteList[clientId] = clientId
}

func (ah authHandler) whiteListOK(clientID string) bool {
	if len(ah.whiteList) == 0 {
		return true
	}

	return ah.whiteList[clientID] == clientID
}

func (ah authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Check for header presence
	authzHeader := r.Header.Get("Authorization")
	if authzHeader == "" {
		log.Println("Missing Authorization header")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized\n"))
		return
	}

	//Header format should be Bearer token
	parts := strings.SplitAfter(authzHeader, "Bearer")
	if len(parts) != 2 {
		log.Println("Unexpected authorization header format - expecting bearer token")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized\n"))
		return
	}

	//Parse the token
	bearerToken := strings.TrimSpace(parts[1])
	token, err := jwt.Parse(bearerToken, roll.GenerateKeyExtractionFunction(ah.secretsRepo))
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized\n"))
		return
	}

	//Make sure the token is valid
	if !token.Valid {
		log.Println("Invalid token presented to service, ", token)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized\n"))
		return
	}

	//Make sure it's no an authcode token
	if roll.IsAuthCode(token) {
		log.Println("Auth code used as access token - access denied")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized\n"))
		return
	}

	//Check against the whitelist
	aud, ok := token.Claims["aud"].(string)
	if !ok {
		log.Println("string aud claim not present in token")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized\n"))
		return
	}
	if !ah.whiteListOK(aud) {
		log.Println("token failed whitelist check:", aud)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized\n"))
		return
	}

	ah.handler.ServeHTTP(w, r)
}
