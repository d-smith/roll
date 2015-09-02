package http

import (
	"github.com/xtraclabs/roll/roll"
	"net/http"
	"errors"
	"html/template"
)


var templates = template.Must(template.ParseFiles("static/authorize.html"))

type authPageContext struct {
	AppName string
	ClientId string
}

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
	required := []string{"client_id","redirect_uri","response_type"}
	for _,v := range required {
		if vals, ok := params[v]; !ok || len(vals) > 1 {
			return false
		}

	}

	return true
}

func validateInputParams(core *roll.Core, r *http.Request) (*roll.Application, error) {
	params := r.URL.Query()

	if params["response_type"][0] != "token" {
		return nil, errors.New("Only token is support for response_type")
	}


	//Client id is application key
	app, err := core.RetrieveApplication(params["client_id"][0])
	if err != nil {
		return nil, err
	}

	if app == nil {
		return nil, errors.New("Invalid client id")
	}


	if app.RedirectUri != params["redirect_uri"][0] {
		return nil, errors.New("redirect_uri does not match registered redirect URIs")
	}

	return app, nil
}

func handleAuthZGet(core *roll.Core, w http.ResponseWriter, r *http.Request) {

	//Check the query params
	if !requiredQueryParamsPresent(r) {
		respondError(w, http.StatusBadRequest, errors.New("Missing required query params or multiple values for single param"))
		return
	}

	//Validate client_id and redirect_uri
	app, err := validateInputParams(core, r)
	if err != nil {
		respondError(w, http.StatusBadRequest,err)
		return
	}

	//Build and return the login page
	pageCtx := &authPageContext{
		AppName:app.ApplicationName,
		ClientId:app.APIKey,
	}

	err = templates.ExecuteTemplate(w, "authorize.html", pageCtx)
	if err != nil {
		respondError(w, http.StatusInternalServerError,err)
		return
	}

}
