package http

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/xtraclabs/roll/authzwrapper"
	"github.com/xtraclabs/roll/roll"
	"net/http"
	"strings"
)

const (
	//DevelopersBaseURI is the base uri for the service.
	DevelopersBaseURI = "/v1/developers"

	//DevelopersURI is for specific resources
	DevelopersURI = DevelopersBaseURI + "/"
)

func handleDevelopersBase(core *roll.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			listDevelopers(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, errors.New("Method not allowed"))
		}
	})
}

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
	subject, scope, err := subjectAndAdminScopeFromRequestCtx(r)
	if err != nil {
		respondError(w, http.StatusInternalServerError, nil)
		return
	}

	devs, err := core.ListDevelopers(subject, scope)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, devs)

}

func subjectAndAdminScopeFromRequestCtx(r *http.Request) (string, bool, error) {
	log.Info("get subject and admin scope from request")
	subject, subjectOk := context.Get(r, authzwrapper.AuthzSubject).(string)
	adminScope, adminOk := context.Get(r, authzwrapper.AuthzAdminScope).(bool)
	if !subjectOk || !adminOk {
		log.Info("error with subject or scope conversion")
		log.Info(subject)
		log.Info(adminScope)
		return "", false, errors.New("System handler misconfiguration")
	}

	return subject, adminScope, nil
}

func retrieveDeveloper(email string, core *roll.Core, w http.ResponseWriter, r *http.Request) {
	if !roll.ValidateEmail(email) {
		respondError(w, http.StatusBadRequest, fmt.Errorf("Invalid email: %s", email))
		return
	}

	subject, scope, err := subjectAndAdminScopeFromRequestCtx(r)
	if err != nil {
		respondError(w, http.StatusInternalServerError, nil)
		return
	}

	dev, err := core.RetrieveDeveloper(email, subject, scope)
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
	email := strings.TrimPrefix(r.RequestURI, DevelopersURI)
	if email == "" {
		respondError(w, http.StatusNotFound, errors.New("Missing resource"))
		return
	}
	retrieveDeveloper(email, core, w, r)
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

	log.Printf("Handling put with payload %v", dev)

	email := strings.TrimPrefix(r.RequestURI, DevelopersURI)

	//If the user included the email inf the body we ignore it. Ignoring it lets us reuse the
	//developer struct for parsing the request, instead of having a projection of the developer
	//structure used to parse the input
	dev.Email = email

	//Extract the subject from the request header based on security mode
	subject, _, err := subjectAndAdminScopeFromRequestCtx(r)
	if err != nil {
		log.Print("Error extracting subject:", err.Error())
		respondError(w, http.StatusInternalServerError, nil)
		return
	}

	//Set the developer id to the subject
	dev.ID = subject

	//Store the developer information
	if err := core.StoreDeveloper(&dev); err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, nil)
}
