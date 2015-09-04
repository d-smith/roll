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
	apiKey := strings.TrimPrefix(r.RequestURI, ApplicationsBaseURI)
	if apiKey == "" {
		respondNotFound(w)
		return
	}

	dev, err := core.RetrieveApplication(apiKey)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err) //TODO - more helpful error messages and better status codes
		return
	}

	if dev == nil {
		respondNotFound(w)
		return
	}

	respondOk(w, dev)
}

func handleApplicationPut(core *roll.Core, w http.ResponseWriter, r *http.Request) {
	var req roll.Application
	if err := parseRequest(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	//Make sure we use the apikey in the resource not any apikey sent in the JSON.
	apiKey := strings.TrimPrefix(r.RequestURI, ApplicationsBaseURI)
	req.APIKey = apiKey

	//Generate a private/public key pair
	private, public, err := secrets.GenerateKeyPair()
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	//Store keys in secrets vault
	err = core.StoreKeysForApp(apiKey, private, public)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	//Store the application definition
	log.Println("storing app def ", req)
	core.StoreApplication(&req)

	respondOk(w, nil)
}
