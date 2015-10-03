package http

import (
	"errors"
	"fmt"
	"github.com/xtraclabs/roll/roll"
	"net/http"
	"strings"
)

const (
	//DevelopersBaseURI is the base uri for the service.
	DevelopersBaseURI = "/v1/developers/"
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

func listDevelopers(core *roll.Core, w http.ResponseWriter, r *http.Request) {
	devs, err := core.ListDevelopers()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, devs)

}

func retrieveDeveloper(email string, core *roll.Core, w http.ResponseWriter, r *http.Request) {
	if !roll.ValidateEmail(email) {
		respondError(w, http.StatusBadRequest, fmt.Errorf("Invalid email: %s", email))
		return
	}

	dev, err := core.RetrieveDeveloper(email)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	if dev == nil {
		respondNotFound(w)
		return
	}

	respondOk(w, dev)
}

func handleDeveloperGet(core *roll.Core, w http.ResponseWriter, r *http.Request) {
	email := strings.TrimPrefix(r.RequestURI, DevelopersBaseURI)
	switch email {
	case "":
		listDevelopers(core, w, r)
	default:
		retrieveDeveloper(email, core, w, r)
	}

}

func handleDeveloperPut(core *roll.Core, w http.ResponseWriter, r *http.Request) {
	var dev roll.Developer
	if err := parseRequest(r, &dev); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	if err := dev.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	email := strings.TrimPrefix(r.RequestURI, DevelopersBaseURI)

	//Ensure the email in the payload is the same as in the resource
	if dev.Email != email {
		respondError(w, http.StatusBadRequest, errors.New("email in body does not match email in request uri"))
		return
	}

	if err := core.StoreDeveloper(&dev); err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, nil)
}
