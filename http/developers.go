package http

import (
	"errors"
	"fmt"
	"github.com/xtraclabs/roll/roll"
	"net/http"
	"strings"
)

const (
	//DevelopersBaseUri is the base uri for the service.
	DevelopersBaseUri = "/v1/developers/"
)

func handleDevelopers(core *roll.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleDeveloperGet(core, w, r)
		case "PUT":
			handleDeveloperPut(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, errors.New("Method not allowed"))
		}
	})
}

func handleDeveloperGet(core *roll.Core, w http.ResponseWriter, r *http.Request) {
	email := strings.TrimPrefix(r.RequestURI, DevelopersBaseUri)
	if !roll.ValidateEmail(email) {
		respondError(w, http.StatusBadRequest, fmt.Errorf("Invalid email: %s", email))
	}

	_, err := core.RetrieveDeveloper(email)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err) //TODO - more helpful error messages and better status codes
		return
	}

	//Todo - marshall the returned body

	respondOk(w, nil)
}

func handleDeveloperPut(core *roll.Core, w http.ResponseWriter, r *http.Request) {
	var req roll.Developer
	if err := parseRequest(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	email := strings.TrimPrefix(r.RequestURI, DevelopersBaseUri)
	if !roll.ValidateEmail(email) {
		respondError(w, http.StatusBadRequest, fmt.Errorf("Invalid email: %s", email))
	}

	core.StoreDeveloper(&req)

	respondOk(w, nil)
}
