package http

import (
	"github.com/xtraclabs/roll/roll"
	"net/http"
	"errors"
)




const (
//AuthorizeBaseUri is the base uri for the service.
	AuthorizeBaseUri = "/oauth2/authorize"
)

func handleAuthorize(core *roll.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleAuthZGet(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, errors.New("Method not allowed"))
		}
	})
}

func requiredQueryParamsPresent(r *http.Request) bool {
	params := r.URL.Query()
	if _, ok := params["client_id"]; !ok {
		return false
	}

	if _, ok := params["redirect_uri"]; !ok {
		return false
	}

	if _, ok := params["response_type"]; !ok {
		return false
	}

	return true
}

func handleAuthZGet(core *roll.Core, w http.ResponseWriter, r *http.Request) {

	//Check the query params
	if !requiredQueryParamsPresent(r) {
		respondError(w, http.StatusMethodNotAllowed, errors.New("Missing required query params"))
		return
	}

	respondOk(w, nil)
}
