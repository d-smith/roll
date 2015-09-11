package http

import (
	"encoding/json"
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/xtraclabs/roll/roll"
	"log"
	"net/http"
)

const (
	//OAuth2TokenBaseURI is the oauth2 token uri
	OAuth2TokenBaseURI = "/oauth2/token"
)

var (
	//ErrInvalidClientDetails is returned when supplied client details don't match those on record
	ErrInvalidClientDetails = errors.New("Invalid application details")

	//ErrRetrievingAppData is generated if the app data assocaited with a client_id (aka api key) cannot be retrieved
	ErrRetrievingAppData = errors.New("Missing or invalid form data")
)

func handleToken(core *roll.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			handleTokenPost(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, errors.New("Method not allowed"))
		}
	})
}

type authCodeContext struct {
	grantType    string
	clientID     string
	clientSecret string
	redirectURI  string
	authCode     string
	username     string
	password     string
}

func (acc *authCodeContext) validate() error {
	switch acc.grantType {
	case "authorization_code":
		return acc.validateAuthCodeGrantType()
	case "password":
		return acc.validatePasswordGrantType()
	default:
		return errors.New("Invalid grant_type")
	}
}

func (acc *authCodeContext) validateAuthCodeGrantType() error {

	if acc.clientID == "" {
		return errors.New("client_id missing from request")
	}

	if acc.clientSecret == "" {
		return errors.New("client_secret missing from request")
	}

	if acc.redirectURI == "" {
		return errors.New("redirect_uri missing from request")
	}

	if acc.authCode == "" {
		return errors.New("code is missing from request")
	}

	return nil
}

func (acc *authCodeContext) validatePasswordGrantType() error {
	if acc.clientID == "" {
		return errors.New("client_id missing from request")
	}

	if acc.clientSecret == "" {
		return errors.New("client_secret missing from request")
	}

	if acc.username == "" {
		return errors.New("username missing from request")
	}

	if acc.password == "" {
		return errors.New("password missing from request")
	}

	return nil
}

type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func validateAndExtractFormParams(r *http.Request) (*authCodeContext, error) {

	acc := &authCodeContext{
		grantType:    r.FormValue("grant_type"),
		clientID:     r.FormValue("client_id"),
		clientSecret: r.FormValue("client_secret"),
		redirectURI:  r.FormValue("redirect_uri"),
		authCode:     r.FormValue("code"),
		username:     r.FormValue("username"),
		password:     r.FormValue("password"),
	}

	return acc, acc.validate()

}

func validateClientDetails(core *roll.Core, ctx *authCodeContext) (*roll.Application, error) {
	app, err := core.RetrieveApplication(ctx.clientID)
	if err != nil {
		log.Println("Error retrieving app data: ", err.Error())
		return nil, ErrRetrievingAppData
	}

	if app == nil {
		return nil, errors.New("Invalid client id")
	}

	if app.APISecret != ctx.clientSecret {
		return nil, ErrInvalidClientDetails
	}

	if ctx.grantType == "authorization_code" && app.RedirectURI != ctx.redirectURI {
		return nil, ErrInvalidClientDetails
	}

	return app, nil
}

func validateCode(secretsRepo roll.SecretsRepo, ctx *authCodeContext) error {
	token, err := jwt.Parse(ctx.authCode, roll.GenerateKeyExtractionFunction(secretsRepo))
	if err != nil {
		return err
	}

	//Make sure the token is valid
	if !token.Valid {
		log.Println("Invalid token presented to service, ", token)
		return errors.New("Invalid authorization code")
	}

	return nil
}

func handleTokenPost(core *roll.Core, w http.ResponseWriter, r *http.Request) {
	//Verify the form params are as expected: grant_type is authorization_code,
	//a code is present, client_id and client_secret are provided, redirect_uri is
	//provided. The content type should be application/x-www-form-urlencoded
	codeContext, err := validateAndExtractFormParams(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	//The grant type was validated above, so at this point we have two grant types to
	//handle: authorization_code and password
	switch codeContext.grantType {
	case "authorization_code":
		handleAuthCodeGrantType(core, w, r, codeContext)
	case "password":
		handlePasswordGrantType(core, w, r, codeContext)
	default:
		//Never say never...
		respondError(w, http.StatusBadRequest, err)
	}
}

func generateAndRespondWithAccessToken(core *roll.Core, app *roll.Application, w http.ResponseWriter) {
	token, err := generateJWT(core, app)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	//Respond with a JSON document included the access_token and a token type of
	//bearer
	at := accessTokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
	}

	atBytes, err := json.Marshal(&at)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(atBytes)
}

func handleAuthCodeGrantType(core *roll.Core, w http.ResponseWriter, r *http.Request, codeContext *authCodeContext) {
	//Verify the client id and secret, plus the redirect_uri by doing a lookup
	//of the app by client id as api key value
	app, err := validateClientDetails(core, codeContext)
	if err != nil {
		switch err {
		case ErrInvalidClientDetails:
			respondError(w, http.StatusBadRequest, ErrInvalidClientDetails)
		default:
			respondError(w, http.StatusInternalServerError, err)
		}

		return
	}

	//Validate the code - it should be a token signed with the users' private key
	if err = validateCode(core.SecretsRepo, codeContext); err != nil {
		respondError(w, http.StatusUnauthorized, err)
		return
	}

	//If everything is cool, generate a JWT access token
	generateAndRespondWithAccessToken(core, app, w)
}

func handlePasswordGrantType(core *roll.Core, w http.ResponseWriter, r *http.Request, codeContext *authCodeContext) {
	//Validate client details
	app, err := validateClientDetails(core, codeContext)
	if err != nil {
		switch err {
		case ErrInvalidClientDetails:
			respondError(w, http.StatusBadRequest, ErrInvalidClientDetails)
		default:
			respondError(w, http.StatusInternalServerError, err)
		}

		return
	}

	//If the client details checkout, authenticate the user credentials
	authenticated, err := authenticateUser(codeContext.username, codeContext.password, app)
	if err != nil {
		log.Println("Error authenticating user: ", err.Error())
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	//If the user credentials don't check out, we're done.
	if !authenticated {
		respondError(w, http.StatusUnauthorized, nil)
		return
	}

	//Create the access token
	generateAndRespondWithAccessToken(core, app, w)

}
