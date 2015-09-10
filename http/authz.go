package http

import (
	"errors"
	"fmt"
	"github.com/xtraclabs/roll/login"
	"github.com/xtraclabs/roll/roll"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var templates = template.Must(template.ParseFiles("../html/authorize.html", "../html/authorize3leg.html"))

type authPageContext struct {
	AppName  string
	ClientID string
}

const (
	//AuthorizeBaseURI is the base uri for obtaining an access token using the implicit flow graph
	AuthorizeBaseURI = "/oauth2/authorize"

	//ValidateBaseURI is the base uri for the authentication callback for the implicit grant flow
	ValidateBaseURI = "/oauth2/validate"
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
	responseType := r.FormValue("response_type")
	if responseType != "token" && responseType != "code" {
		return nil, errors.New("response_type must be code or token")
	}

	//Client id is application key
	clientID := r.FormValue("client_id")
	app, err := core.RetrieveApplication(clientID)
	if err != nil {
		return nil, err
	}

	if app == nil {
		return nil, errors.New("Invalid client id")
	}

	redirectURI := r.FormValue("redirect_uri")
	if app.RedirectURI != redirectURI {
		return nil, errors.New("redirect_uri does not match registered redirect URIs")
	}

	return app, nil
}

func executeAuthTemplate(w http.ResponseWriter, r *http.Request, pageCtx *authPageContext) error {
	responseType := r.FormValue("response_type")
	var authPage string

	switch responseType {
	case "token":
		authPage = "authorize.html"
	case "code":
		authPage = "authorize3leg.html"
	default:
		authPage = "error"
	}

	if authPage == "error" {
		return errors.New("Unable to build authorization page for response_type " + responseType)
	}

	return templates.ExecuteTemplate(w, authPage, pageCtx)

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
		ClientID: app.APIKey,
	}

	err = executeAuthTemplate(w, r, pageCtx)
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

func lookupApplicationFromFormClientID(core *roll.Core, r *http.Request) (*roll.Application, error) {
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

func buildDeniedRedirectURL(app *roll.Application) string {
	return fmt.Sprintf("%s#error=access_denied", app.RedirectURI)
}

func buildRedirectURL(core *roll.Core, w http.ResponseWriter, responseType string, app *roll.Application) (string, error) {
	log.Println("build redirect, app ctx:", app.RedirectURI)

	var redirectURL string
	switch responseType {
	case "token":
		//Create signed token
		token, err := generateJWT(core, app)
		if err != nil {
			return "", err
		}
		redirectURL = fmt.Sprintf("%s#access_token=%s&token_type=Bearer", app.RedirectURI, token)
	case "code":
		token, err := generateSignedCode(core, app)
		if err != nil {
			return "", err
		}
		redirectURL = fmt.Sprintf("%s?code=%s", app.RedirectURI, token)
	default:
		panic(errors.New("unexpected response type in buildRedirectURL: " + responseType))
	}
	return redirectURL, nil
}

func generateJWT(core *roll.Core, app *roll.Application) (string, error) {
	privateKey, err := core.RetrievePrivateKeyForApp(app.APIKey)
	if err != nil {
		return "", err
	}

	token, err := roll.GenerateToken(app, privateKey)
	return token, err
}

func generateSignedCode(core *roll.Core, app *roll.Application) (string, error) {
	privateKey, err := core.RetrievePrivateKeyForApp(app.APIKey)
	if err != nil {
		return "", err
	}

	token, err := roll.GenerateCode(app, privateKey)
	return token, err
}

func getResponseType(r *http.Request) (string, error) {
	if len(r.Form["response_type"]) != 1 {
		return "", errors.New("Expected single response_type param as part of query params")
	}

	responseType := r.Form["response_type"][0]
	if responseType != "token" && responseType != "code" {
		return "", errors.New("valid vaule for response_type are token and code")
	}

	return responseType, nil

}

func authenticateUser(username, password string, app *roll.Application) (bool, error) {
	//Convert login provider to a URL
	loginURL, err := url.Parse(app.LoginProvider)
	if err != nil {
		return false, err
	}

	//Grab the login kit
	kit := login.GetLoginKit(loginURL.Scheme)
	if kit == nil {
		return false, errors.New("No login kit for login provider " + loginURL.Scheme)
	}

	//Form the request and endpoint
	loginRequest := kit.RequestBuilder(username, password)
	endpoint := kit.EndpointBuilder(loginURL.Host)

	//Send it
	client := &http.Client{}
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(loginRequest))
	if err != nil {
		return false, err
	}

	req.Header.Add("SOAPAction", "\"\"")
	req.Header.Add("Content-Type", "text/xml")

	log.Println(fmt.Sprintf("%v\n%s", req, loginRequest))

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	log.Println(fmt.Sprintf("%v", resp))

	return resp.StatusCode == http.StatusOK, nil

}

func handleAuthZValidate(core *roll.Core, w http.ResponseWriter, r *http.Request) {

	//Parse request form
	err := r.ParseForm()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	log.Println(r.Form)

	responseType, err := getResponseType(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	//Lookup the client based on the hidden input field
	app, err := lookupApplicationFromFormClientID(core, r)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
	}

	//Check if user denied authorization
	if denied(r) {
		redirectURL := buildDeniedRedirectURL(app)
		http.Redirect(w, r, redirectURL, http.StatusFound)
		return
	}

	authenticated, err := authenticateUser(r.Form["username"][0], r.Form["password"][0], app)
	if err != nil {
		log.Println("Error authenticating user: ", err.Error())
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	//Was the authentication successful?
	if !authenticated {
		redirectURL := buildDeniedRedirectURL(app)
		http.Redirect(w, r, redirectURL, http.StatusFound)
		return
	}

	//Build redirect url with embedded token or code
	redirectURL, err := buildRedirectURL(core, w, responseType, app)
	if err != nil {
		log.Println("Error generating redirect url: ", err.Error())
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	//Redirect the user to the new URL
	http.Redirect(w, r, redirectURL, http.StatusFound)

}
