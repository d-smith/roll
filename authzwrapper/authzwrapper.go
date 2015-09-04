package authzwrapper

import (
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/xtraclabs/roll/repos"
	"github.com/xtraclabs/roll/roll"
	"log"
	"net/http"
	"strings"
)

//AuthHandler is a wrapper type, embedding the wrapped handler and the secrets repo needed to
//look up the key for validating JWT bearer tokens
type AuthHandler struct {
	handler     http.Handler
	secretsRepo roll.SecretsRepo
}

//Wrap takes a handler and decorates it with JWT bearer token validation.
func Wrap(h http.Handler) *AuthHandler {
	return &AuthHandler{
		handler:     h,
		secretsRepo: repos.NewVaultSecretsRepo(),
	}
}

func (ah *AuthHandler) makeKeyExtractionFunction() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		//Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		//The api key is carried in aud
		apiKey := token.Claims["aud"]
		if apiKey == "" {
			return nil, errors.New("api key not found in aud claim")
		}

		//get the public key from the vault
		keystring, err := ah.secretsRepo.RetrievePublicKeyForApp(apiKey.(string))
		if err != nil {
			return nil, err
		}

		//Parse the keystring
		return jwt.ParseRSAPublicKeyFromPEM([]byte(keystring))
	}
}

func (ah *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	token, err := jwt.Parse(bearerToken, ah.makeKeyExtractionFunction())
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

	ah.handler.ServeHTTP(w, r)
}
