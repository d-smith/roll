package http

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/rollsecrets/secrets"
	rolltoken "github.com/xtraclabs/rollsecrets/token"
	"net/http"
	"strings"
)

const (
	//OAuth2TokenBaseURI is the oauth2 token uri
	OAuth2TokenBaseURI = "/oauth2/token"
)

var (
	//ErrInvalidClientDetails is returned when supplied client details don't match those on record
	ErrInvalidClientDetails = errors.New("Invalid application details")

	//ErrRetrievingAppData is generated if the app data associated with a client_id (aka api key) cannot be retrieved
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
	assertion    string
	scope        string
}

func (acc *authCodeContext) validate() error {
	switch acc.grantType {
	case "authorization_code":
		return acc.validateAuthCodeGrantType()
	case "password":
		return acc.validatePasswordGrantType()
	case "urn:ietf:params:oauth:grant-type:jwt-bearer":
		return acc.validateJWTGrantType()
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

func (acc *authCodeContext) validateJWTGrantType() error {
	if acc.assertion == "" {
		return errors.New("assertion missing from request")
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
		assertion:    r.FormValue("assertion"),
		scope:        r.FormValue("scope"),
	}

	return acc, acc.validate()

}

func subjectFromBearerToken(core *roll.Core, r *http.Request) (string, error) {
	//Check for header presence
	authzHeader := r.Header.Get("Authorization")
	if authzHeader == "" {
		return "", errors.New("Authorization header missing from request")
	}

	//Header format should be Bearer token
	parts := strings.SplitAfter(authzHeader, "Bearer")
	if len(parts) != 2 {
		return "", errors.New("Unexpected authorization header format - expecting bearer token")
	}

	//Parse the token
	bearerToken := strings.TrimSpace(parts[1])
	token, err := jwt.Parse(bearerToken, rolltoken.GenerateKeyExtractionFunction(core.SecretsRepo))
	if err != nil {
		return "", err
	}

	//Grab the subject from the claims
	subject, ok := token.Claims["sub"].(string)
	if !ok {
		return "", errors.New("problem with subject claim")
	}

	//Is the subject something other than an empty string?
	if subject == "" {
		return "", errors.New("empty subject claim")
	}

	return subject, nil
}

func lookupApplication(core *roll.Core, clientID string) (*roll.Application, error) {
	app, err := core.SystemRetrieveApplication(clientID)
	if err != nil {
		log.Info("Error retrieving app data: ", err.Error())
		return nil, ErrRetrievingAppData
	}

	if app == nil {
		log.Info("Invalid client id: ", clientID)
		return nil, errors.New("Invalid client id")
	}

	return app, nil
}

func lookupApplicatioByAudience(core *roll.Core, audience string) (*roll.Application, error) {
	app, err := core.SystemRetrieveApplicationByJWTFlowAudience(audience)
	if err != nil {
		log.Info("Error retrieving app data: ", err.Error())
		return nil, ErrRetrievingAppData
	}

	if app == nil {
		log.Info("invalid client id")
		return nil, errors.New("Invalid client id")
	}

	return app, nil
}

func validateClientDetails(core *roll.Core, ctx *authCodeContext) (*roll.Application, error) {
	app, err := lookupApplication(core, ctx.clientID)
	if err != nil {
		log.Info("error looking up application")
		return nil, err
	}

	if app.ClientSecret != ctx.clientSecret {
		log.Info("error validating client secret")
		log.Info("secret from db: ", app.ClientSecret)
		log.Info("secret from context: ", ctx.clientSecret)

		return nil, ErrInvalidClientDetails
	}

	if ctx.grantType == "authorization_code" && app.RedirectURI != ctx.redirectURI {
		log.Info("error validating registered redirect URI")
		return nil, ErrInvalidClientDetails
	}

	return app, nil
}

func validateAndReturnCodeToken(secretsRepo secrets.SecretsRepo, ctx *authCodeContext, clientID string) (*jwt.Token, error) {
	token, err := jwt.Parse(ctx.authCode, rolltoken.GenerateKeyExtractionFunction(secretsRepo))
	if err != nil {
		return nil, err
	}

	//Make sure the token is valid
	if !token.Valid {
		log.Info("Invalid token presented to service: ", token)
		return nil, errors.New("Invalid authorization code")
	}

	//make sure the client_id used to validate the token matches the token aud claim
	if clientID != token.Claims["aud"] {
		return nil, errors.New("Token not associated with client id")
	}

	return token, nil
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
	case "urn:ietf:params:oauth:grant-type:jwt-bearer":
		handleJWTGrantType(core, w, r, codeContext)
	default:
		//Never say never...
		respondError(w, http.StatusBadRequest, err)
	}
}

func generateAndRespondWithAccessToken(core *roll.Core, subject, scope string, app *roll.Application, w http.ResponseWriter) {
	token, err := generateJWT(subject, scope, core, app)
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
	token, err := validateAndReturnCodeToken(core.SecretsRepo, codeContext, r.FormValue("client_id"))
	if err != nil {
		respondError(w, http.StatusUnauthorized, err)
		return
	}

	scope, ok := token.Claims["scope"].(string)
	if !ok {
		respondError(w, http.StatusBadRequest, errors.New("problem with token scope"))
		return
	}

	subject, ok := token.Claims["sub"].(string)
	if !ok {
		respondError(w, http.StatusBadRequest, errors.New("problem with token subject"))
		return
	}

	//If everything is cool, generate a JWT access token
	generateAndRespondWithAccessToken(core, subject, scope, app, w)
}

func handlePasswordGrantType(core *roll.Core, w http.ResponseWriter, r *http.Request, codeContext *authCodeContext) {
	//Validate client details
	log.Info("Handle password grant type")
	log.Info("Validate client details")
	app, err := validateClientDetails(core, codeContext)
	if err != nil {
		log.Info(err.Error())
		switch err {
		case ErrInvalidClientDetails:
			respondError(w, http.StatusBadRequest, ErrInvalidClientDetails)
		default:
			respondError(w, http.StatusInternalServerError, err)
		}

		return
	}

	//If the client details checkout, authenticate the user credentials
	log.Info("authenticate user credentials")
	authenticated, err := authenticateUser(codeContext.username, codeContext.password, app)
	if err != nil {
		log.Info("Error authenticating user: ", err.Error())
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	//If the user credentials don't check out, we're done.
	if !authenticated {
		respondError(w, http.StatusUnauthorized, nil)
		return
	}

	//If a scope is present, validate it.
	log.Info("validate scope")
	valid, err := validateScopes(core, r)
	if err != nil {
		log.Info("error validating scope: ", err.Error())
		respondError(w, http.StatusInternalServerError, nil)
		return
	}

	if !valid {
		log.Info("scope is invalid")
		respondError(w, http.StatusUnauthorized, nil)
		return
	}

	//Create the access token
	generateAndRespondWithAccessToken(core, codeContext.username, codeContext.scope, app, w)

}

func filterUnsupportedClaims(scope string) string {
	if scope == "" {
		return scope
	}

	var supportedScope string

	scopeParts := strings.Fields(scope)
	for _, sp := range scopeParts {
		//currently only admin is supported
		if sp == "admin" {
			supportedScope = sp
		}
	}

	return supportedScope
}

func handleJWTGrantType(core *roll.Core, w http.ResponseWriter, r *http.Request, codeContext *authCodeContext) {
	log.Info("handleJWTGrantType")

	//First step is to verify the token signature
	log.Info("verify token signature")
	token, err := jwt.Parse(codeContext.assertion, generateKeyExtractionFunctionForJTWFlow(core.ApplicationRepo))
	if err != nil {
		log.Info(err.Error())
		respondError(w, http.StatusUnauthorized, err)
		return
	}

	//Grab the app definition based on iss carries the api key/client_id
	log.Info("look up application definition")
	app, err := lookupApplicatioByAudience(core, token.Claims["aud"].(string))
	if err != nil {
		switch err {
		case ErrInvalidClientDetails:
			respondError(w, http.StatusBadRequest, ErrInvalidClientDetails)
		default:
			respondError(w, http.StatusInternalServerError, err)
		}

		return
	}

	//Pass the identity
	subject, ok := token.Claims["sub"].(string)
	if !ok {
		respondError(w, http.StatusBadRequest, errors.New("sub claim is not a string"))
		return
	}

	//Make sure the claim conveys a sub
	if subject == "" {
		respondError(w, http.StatusBadRequest, errors.New("JWT missing sub claim"))
		return
	}

	//Include scope
	scope, ok := token.Claims["scope"].(string)
	if !ok {
		scope = ""
	}

	//Now we can generate a token since we had the app needed to form the token
	log.Info("generate token")

	//TODO - extract and validate scope
	generateAndRespondWithAccessToken(core, subject, filterUnsupportedClaims(scope), app, w)

}

func generateKeyExtractionFunctionForJTWFlow(applicationRepo roll.ApplicationRepo) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		//Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		//The aud claim conveys the intended application the token is used to gain access to.
		clientID := token.Claims["aud"]
		if clientID == nil {
			return nil, errors.New("Foreign token does not include aud claim")
		}

		//Look up the application
		app, err := applicationRepo.SystemRetrieveApplicationByJWTFlowAudience(clientID.(string))
		if err != nil {
			log.Info("Error looking up app for ", clientID, " ", err.Error())
			return nil, err
		}

		if app == nil {
			log.Info("No app definition associated with audience found: ", clientID.(string))
			return nil, errors.New("No app definition associated with aud found")
		}

		//We also check that the token was issued by the entity registered with the application
		issuer := token.Claims["iss"]
		if issuer == nil || issuer != app.JWTFlowIssuer {
			return nil, errors.New("Foreign token issuer not known")
		}

		//Grab the public key from the app definition
		keystring := app.JWTFlowPublicKey

		log.Info("validating with '", keystring, "'")

		//Parse the keystring
		return jwt.ParseRSAPublicKeyFromPEM([]byte(keystring))
	}
}
