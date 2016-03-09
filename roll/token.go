package roll

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/xtraclabs/rollsecrets/secrets"
	"strings"
	"time"
)

//XtAuthCodeScope defines a constant to use as a scope claim to limit use of auth code for the 3 legged
//OAuth2 flow
const XtAuthCodeScope = "xtAuthCode"

var idGenerator IdGenerator = UUIDIdGenerator{}

func scopeStringWithoutAuthcodeScope(scope string) string {
	if scope == "" {
		return scope
	}

	scopeParts := strings.Fields(scope)
	var finalScope string
	for _, sp := range scopeParts {
		if sp != XtAuthCodeScope {
			if len(finalScope) != 0 {
				finalScope += " "
			}
			finalScope += sp
		}
	}

	return finalScope
}

//GenerateToken generates a signed JWT for an application using the
//provided private key which is assocaited with the application.
func GenerateToken(subject, scope string, app *Application, privateKey string) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))

	jti, err := idGenerator.GenerateID()
	if err != nil {
		return "", err
	}

	t.Claims["sub"] = subject
	t.Claims["aud"] = app.ClientID
	t.Claims["iat"] = int64(time.Now().Unix())
	t.Claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	t.Claims["jti"] = jti
	t.Claims["application"] = app.ApplicationName

	//Add scope to claim, removing our reserved auth code scope
	t.Claims["scope"] = scopeStringWithoutAuthcodeScope(scope)

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return "", nil
	}

	tokenString, err := t.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

//GenerateCode generates a code for the web server callback in a 3-legged Oauth2 flow. We create
//these as signed tokens so we can see a) if its one of ours and b) if its still valid.
func GenerateCode(subject, scope string, app *Application, privateKey string) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))

	jti, err := idGenerator.GenerateID()
	if err != nil {
		return "", err
	}

	t.Claims["sub"] = subject
	t.Claims["aud"] = app.ClientID
	t.Claims["jti"] = jti
	t.Claims["exp"] = time.Now().Add(30 * time.Second).Unix()

	codeScope := XtAuthCodeScope
	if scope != "" {
		codeScope += " " + scope
	}
	t.Claims["scope"] = codeScope

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return "", nil
	}

	tokenString, err := t.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

//GenerateKeyExtractionFunction generates a key extraction function for verifying the signature
//associated with a JWT
func GenerateKeyExtractionFunction(secretsRepo secrets.SecretsRepo) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		//Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		//The api key is carried in aud
		clientID := token.Claims["aud"]
		if clientID == "" {
			return nil, errors.New("api key not found in aud claim")
		}

		//get the public key from the vault
		keystring, err := secretsRepo.RetrievePublicKeyForApp(clientID.(string))
		if err != nil {
			return nil, err
		}

		//Parse the keystring
		return jwt.ParseRSAPublicKeyFromPEM([]byte(keystring))
	}
}

func GenerateKeyExtractionFunctionForJTWFlow(applicationRepo ApplicationRepo) jwt.Keyfunc {
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

func IsAuthCode(token *jwt.Token) bool {
	if !token.Valid {
		panic(errors.New("Attempting to test invalid token for auth code claim"))
	}

	return token.Claims["scope"] == XtAuthCodeScope
}
