package http

import (
	"errors"
	"github.com/xtraclabs/roll/repos"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/roll/secrets"
	"log"
	"net/http"
	"strings"
)

const (
	//ApplicationsBaseURI is the base uri for the service.
	ApplicationsBaseURI = "/v1/applications"

	//ApplicationsURI is for specific resources
	ApplicationsURI = ApplicationsBaseURI + "/"
)

type ApplicationCreatedResponse struct {
	ClientID string `json:"client_id"`
}

func handleApplicationsBase(core *roll.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			handleApplicationPost(core, w, r)
		case "GET":
			listApplications(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, errors.New("Method not allowed"))
		}
	})
}
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

func retrieveApplication(clientID string, core *roll.Core, w http.ResponseWriter, r *http.Request) {
	log.Println("ret appl called", clientID)
	if clientID == "" {
		respondError(w, http.StatusBadRequest, errors.New("Resource not specified"))
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

func listApplications(core *roll.Core, w http.ResponseWriter, r *http.Request) {
	subject, scope, err := subjectAndAdminScopeFromRequestCtx(r)
	if err != nil {
		respondError(w, http.StatusInternalServerError, nil)
		return
	}

	apps, err := core.ListApplications(subject, scope)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, apps)
}

func handleApplicationGet(core *roll.Core, w http.ResponseWriter, r *http.Request) {
	log.Println("handleApplicationGet called")
	clientID := strings.TrimPrefix(r.RequestURI, ApplicationsURI)
	retrieveApplication(clientID, core, w, r)
}

func handleApplicationPost(core *roll.Core, w http.ResponseWriter, r *http.Request) {
	var app roll.Application
	if err := parseRequest(r, &app); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	//Assign a client ID
	id, err := core.GenerateID()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	app.ClientID = id

	//Validate the content
	if err := app.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	//Extract the subject from the request header based on security mode
	subject, err := subjectFromAuthHeader(core, r)
	if err != nil {
		log.Print("Error extracting subject:", err.Error())
		respondError(w, http.StatusInternalServerError, nil)
		return
	}

	app.DeveloperID = subject

	//Store the application definition
	log.Println("storing app def ", app)
	err = core.CreateApplication(&app)
	if err != nil {
		switch err.(type) {
		case *repos.DuplicateAppdefError:
			respondError(w, http.StatusConflict, err)
		default:
			respondError(w, http.StatusInternalServerError, err)
		}

		return
	}

	//Generate a private/public key pair
	private, public, err := secrets.GenerateKeyPair()
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	//Store keys in secrets vault
	err = core.StoreKeysForApp(id, private, public)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	//Return the client id
	clientID := ApplicationCreatedResponse{ClientID: id}

	respondOk(w, clientID)

}

func handleApplicationPut(core *roll.Core, w http.ResponseWriter, r *http.Request) {
	var app roll.Application
	if err := parseRequest(r, &app); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	//Make sure we use the clientID in the resource not any clientID sent in the JSON.
	clientID := strings.TrimPrefix(r.RequestURI, ApplicationsURI)
	if clientID == "" {
		respondError(w, http.StatusBadRequest, nil)
		return
	}

	app.ClientID = clientID

	//Validate the content
	if err := app.Validate(); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	//Retrieve the app definition to update
	storedApp, err := core.RetrieveApplication(clientID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	if storedApp == nil {
		respondError(w, http.StatusNotFound, nil)
		return
	}

	//Copy over the potential updates
	storedApp.ApplicationName = app.ApplicationName
	storedApp.DeveloperEmail = app.DeveloperEmail
	storedApp.LoginProvider = app.LoginProvider
	storedApp.RedirectURI = app.RedirectURI
	storedApp.DeveloperID = app.DeveloperID

	//Store the application definition
	log.Println("updating app def ", app)
	err = core.UpdateApplication(&app)
	if err != nil {
		log.Println("Error updating definition: ", err.Error())
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, nil)
}
