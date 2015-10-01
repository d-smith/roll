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
}

//Wrap takes a handler and decorates it with JWT bearer token validation.
func Wrap(secretsRepo roll.SecretsRepo, h http.Handler) http.Handler {
	return &authHandler{
		handler:     h,
		secretsRepo: secretsRepo,
	}
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

	ah.handler.ServeHTTP(w, r)
}
