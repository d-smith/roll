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

	//ErrTokenParsing is generated if the auth code in the form of a JWT cannot be parsed
	ErrTokenParsing         = errors.New("Invalid authorization code")

	//ErrRetrievingAppData is generated if the app data assocaited with a client_id (aka api key) cannot be retrieved
	ErrRetrievingAppData    = errors.New("Missing or invalid form data")
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
	clientID string
	clientSecret string
	redirectURI  string
	authCode     string
}

type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func validateAndExtractFormParams(r *http.Request) (*authCodeContext, error) {
	grantType := r.FormValue("grant_type")
	if grantType != "authorization_code" {
		return nil, errors.New("Invalid grant_type")
	}

	clientID := r.FormValue("client_id")
	if clientID == "" {
		return nil, errors.New("client_id missing from request")
	}

	clientSecret := r.FormValue("client_secret")
	if clientSecret == "" {
		return nil, errors.New("client_secret missing from request")
	}

	redirectURI := r.FormValue("redirect_uri")
	if redirectURI == "" {
		return nil, errors.New("redirect_uri missing from request")
	}

	authCode := r.FormValue("code")
	if authCode == "" {
		return nil, errors.New("code is missing from request")
	}

	return &authCodeContext{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
		authCode:     authCode,
	}, nil

}

func validateClientDetails(core *roll.Core, ctx *authCodeContext) (*roll.Application, error) {
	app, err := core.RetrieveApplication(ctx.clientID)
	if err != nil {
		log.Println("Error retrieving app data: ", err.Error())
		return nil, ErrRetrievingAppData
	}

	if app.APISecret != ctx.clientSecret {
		return nil, ErrInvalidClientDetails
	}

	if app.RedirectURI != ctx.redirectURI {
		return nil, ErrInvalidClientDetails
	}

	return app, nil

}

func validateCode(secretsRepo roll.SecretsRepo, ctx *authCodeContext) error {
	token, err := jwt.Parse(ctx.authCode, roll.GenerateKeyExtractionFunction(secretsRepo))
	if err != nil {
		return ErrTokenParsing
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
		switch err {
		case ErrTokenParsing:
			respondError(w, http.StatusInternalServerError, err)
		default:
			respondError(w, http.StatusBadRequest, err)
		}

		return
	}

	//If everything is cool, generate a JWT access token
	token, err := generateJWT(core, app)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
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
