package http

import (
	"errors"
	"fmt"
	"github.com/xtraclabs/roll/roll"
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseFiles("../html/authorize.html"))

type authPageContext struct {
	AppName  string
	ClientId string
}

const (
	AuthorizeBaseUri = "/oauth2/authorize"
	ValidateBaseUri  = "/oauth2/validate"
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
	required := []string{"client_id", "redirect_uri", "response_type"}
	for _, v := range required {
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
		respondError(w, http.StatusBadRequest, err)
		return
	}

	//Build and return the login page
	pageCtx := &authPageContext{
		AppName:  app.ApplicationName,
		ClientId: app.APIKey,
	}

	err = templates.ExecuteTemplate(w, "authorize.html", pageCtx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

}

func handleValidate(core *roll.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			handleAuthZValidate(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, errors.New("Method not allowed"))
		}
	})
}

func denied(r *http.Request) bool {
	authstate := r.Form["authorize"][0]
	return authstate != "allow"
}

func lookupApplicationFromFormClientId(core *roll.Core, r *http.Request) (*roll.Application, error) {
	app, err := core.RetrieveApplication(r.Form["client_id"][0])
	if err != nil {
		return nil, err
	}

	if app == nil {
		return nil, errors.New("Invalid client id")
	}

	//TODO - separate error handling for internal error and invalid client id

	return app, nil
}

func buildDeniedRedirectUrl(app *roll.Application) string {
	return fmt.Sprintf("%s#error=access_denied", app.RedirectUri)
}

func buildRedirectUrl(token string, app *roll.Application) string {
	return fmt.Sprintf("%s#access_token=%s&token_type=Bearer", app.RedirectUri, token)
}

func generateJWT(core *roll.Core, app *roll.Application) (string, error) {
	privateKey, err := core.RetrievePrivateKeyForApp(app.APIKey)
	if err != nil {
		return "", err
	}

	token, err := roll.GenerateToken(app, privateKey)
	return token, err
}

func handleAuthZValidate(core *roll.Core, w http.ResponseWriter, r *http.Request) {

	//Parse request form
	err := r.ParseForm()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	//Lookup the client based on the hidden input field
	app, err := lookupApplicationFromFormClientId(core, r)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
	}

	//Check if user denied authorization
	if denied(r) {
		redirectURL := buildDeniedRedirectUrl(app)
		http.Redirect(w, r, redirectURL, http.StatusFound)
		return
	}

	//Fake out the authenticaiton
	log.Println("WARNING - CURRENTLY NO LOOKUP OF LOGIN ENDPOINT AND AUTHENTICATION CALL")

	//Fake out token creation
	token, err := generateJWT(core, app)
	if err != nil {
		log.Println("Error generating token: ", err.Error())
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	//Build redirect url
	redirectURL := buildRedirectUrl(token, app)

	//Redirect the user to the new URL
	http.Redirect(w, r, redirectURL, http.StatusFound)

}
