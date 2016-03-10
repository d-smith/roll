package authzwrapper

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/rollsecrets/secrets"
	"github.com/xtraclabs/rollauthz/rollauthz"
	"net/http"
)

type key int

//Use to get the request context, from which the subject can be extracted (in both secure and unsecure modes)
const AuthzSubject key = 0
const AuthzAdminScope key = 1

type authHandler struct {
	handler     http.Handler
	rollAuthZ 	*rollauthz.RollAuthZ
	adminRepo   roll.AdminRepo
	whiteList   map[string]string
}

//Wrap takes a handler and decorates it with JWT bearer token validation.
func Wrap(secretsRepo secrets.SecretsRepo, adminRepo roll.AdminRepo, whitelistedClientIDs []string, h http.Handler) http.Handler {
	wl := make(map[string]string)
	for _, cid := range whitelistedClientIDs {
		wl[cid] = cid
	}

	return &authHandler{
		handler:     h,
		rollAuthZ: &rollauthz.RollAuthZ{secretsRepo},
		adminRepo:   adminRepo,
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
		log.Info("Missing Authorization header")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized\n"))
		return
	}

	claims, err := ah.rollAuthZ.ValidateAccessToken(authzHeader)
	if err != nil {
		log.Info(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized\n"))
	}

	//Check against the whitelist
	aud, ok := claims["aud"].(string)
	if !ok {
		log.Info("aud claim not present in token")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized\n"))
		return
	}

	if !ah.whiteListOK(aud) {
		log.Info("token failed whitelist check: ", aud)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized\n"))
		return
	}


	sub, ok := claims["sub"].(string)
	if !ok {
		log.Info("Unable to extract sub from token claims")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized\n"))
		return
	}

	context.Set(r, AuthzAdminScope, false)
	scope, ok := claims["scope"].(string)
	if ok && scope == "admin" {
		admin, err := ah.adminRepo.IsAdmin(sub)
		if err != nil {
			log.Info("error making admin scope determination: ", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized\n"))
			return
		}

		context.Set(r, AuthzAdminScope, admin)
	}

	context.Set(r, AuthzSubject, sub)
	ah.handler.ServeHTTP(w, r)
	context.Clear(r)
}
