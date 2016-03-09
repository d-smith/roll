package http

import (
	"encoding/json"
	"errors"
	log "github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/xtraclabs/roll/roll"
	rolltoken "github.com/xtraclabs/rollsecrets/token"
	"net/http"
)

const (
	//TokenInfoURI is the base uri for the token validation service.
	TokenInfoURI = "/oauth2/tokeninfo"
)

type tokenInfo struct {
	Audience string `json:"audience"`
}

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
		log.Info("missing access token")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Parse and validate the token
	token, err := jwt.Parse(tokenString, rolltoken.GenerateKeyExtractionFunction(core.SecretsRepo))
	if err != nil {
		log.Info(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Panic is there's no aud claim
	audience := token.Claims["aud"]
	if audience == "" {
		panic(errors.New("No aud claim in token"))
	}

	//Return the token info
	tokenInfo := &tokenInfo{
		Audience: audience.(string),
	}

	bytes, err := json.Marshal(&tokenInfo)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	w.Write(bytes)
}
