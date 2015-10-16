package roll

import (
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

//XtAuthCodeScope defines a constant to use as a scope claim to limit use of auth code for the 3 legged
//OAuth2 flow
const XtAuthCodeScope = "xtAuthCode"

var idGenerator IdGenerator = UUIDIdGenerator{}

//GenerateToken generates a signed JWT for an application using the
//provided private key which is assocaited with the application.
func GenerateToken(app *Application, privateKey string) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))

	jti, err := idGenerator.GenerateID()
	if err != nil {
		return "", err
	}

	t.Claims["aud"] = app.ClientID
	t.Claims["iat"] = int64(time.Now().Unix())
	t.Claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	t.Claims["jti"] = jti
	t.Claims["application"] = app.ApplicationName

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
func GenerateCode(app *Application, privateKey string) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))

	jti, err := idGenerator.GenerateID()
	if err != nil {
		return "", err
	}

	t.Claims["aud"] = app.ClientID
	t.Claims["jti"] = jti
	t.Claims["exp"] = time.Now().Add(30 * time.Second).Unix()
	t.Claims["scope"] = XtAuthCodeScope

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
func GenerateKeyExtractionFunction(secretsRepo SecretsRepo) jwt.Keyfunc {
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

		//The api key is carried in iss
		clientID := token.Claims["iss"]
		if clientID == "" {
			return nil, errors.New("api key not found in iss claim")
		}

		//Look up the application
		app, err := applicationRepo.RetrieveApplication(clientID.(string))
		if err != nil {
			return nil, err
		}

		if app == nil {
			return nil, errors.New("No app definition associated with iss found")
		}

		//Grab the public key from the app definition
		keystring := app.JWTFlowPublicKey

		log.Println("validating with '", keystring, "'")

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
