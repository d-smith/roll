package http

import (
	"errors"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/roll/secrets"
	"log"
	"net/http"
	"strings"
)

const (
	//ApplicationsBaseURI is the base uri for the service.
	ApplicationsBaseURI = "/v1/applications/"
)

func handleApplications(core *roll.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleApplicationGet(core, w, r)
		case "PUT":
			handleApplicationPut(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, errors.New("Method not allowed"))
		}
	})
}

func handleApplicationGet(core *roll.Core, w http.ResponseWriter, r *http.Request) {
	clientID := strings.TrimPrefix(r.RequestURI, ApplicationsBaseURI)
	if clientID == "" {
		respondNotFound(w)
		return
	}

	app, err := core.RetrieveApplication(clientID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	if app == nil {
		respondNotFound(w)
		return
	}

	respondOk(w, app)
}

func handleApplicationPut(core *roll.Core, w http.ResponseWriter, r *http.Request) {
	var app roll.Application
	if err := parseRequest(r, &app); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	//Make sure we use the clientID in the resource not any clientID sent in the JSON.
	clientID := strings.TrimPrefix(r.RequestURI, ApplicationsBaseURI)
	app.ClientID = clientID

	//Validate the content
	if err := app.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	//Generate a private/public key pair
	private, public, err := secrets.GenerateKeyPair()
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	//Store keys in secrets vault
	err = core.StoreKeysForApp(clientID, private, public)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	//Store the application definition
	log.Println("storing app def ", app)
	err = core.StoreApplication(&app)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, nil)
}
