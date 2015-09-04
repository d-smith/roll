package authzwrapper

import (
	"net/http"
	jwt "github.com/dgrijalva/jwt-go"
	"errors"
	"github.com/xtraclabs/roll/repos"
	"github.com/xtraclabs/roll/roll"
	"fmt"
	"log"
)

type AuthHandler struct {
	handler http.Handler
	secretsRepo roll.SecretsRepo
}

func Wrap(h http.Handler) *AuthHandler {
	return &AuthHandler {
		handler:h,
		secretsRepo: repos.NewVaultSecretsRepo(),
	}
}



func (ah *AuthHandler)  makeKeyExtractionFunction() jwt.Keyfunc {
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
		keystring, err :=  ah.secretsRepo.RetrievePublicKeyForApp(apiKey.(string))
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
	if authzHeader ==  "" {
		log.Println("Missing Authorization header")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized\n"))
		return
	}

	//Parse the token
	token, err := jwt.Parse(authzHeader, ah.makeKeyExtractionFunction())
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

	ah.handler.ServeHTTP(w,r)
}