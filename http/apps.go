package http

import (
	"github.com/xtraclabs/roll/roll"
	"net/http"
	"errors"
	"strings"
	"log"
)

const (
	//ApplicationsBaseUri is the base uri for the service.
	ApplicationsBaseUri = "/v1/applications/"
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
	apiKey := strings.TrimPrefix(r.RequestURI, ApplicationsBaseUri)

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

	log.Println("storing app def ", req)
	core.StoreApplication(&req)

	respondOk(w, nil)
}