package roll

import (
	"errors"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/nu7hatch/gouuid"
	"time"
)

//GenerateToken generates a signed JWT for an application using the
//provided private key which is assocaited with the application.
func GenerateToken(app *Application, privateKey string) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))

	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	t.Claims["aud"] = app.APIKey
	t.Claims["iat"] = int64(time.Now().Unix())
	t.Claims["jti"] = u.String()
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

	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	t.Claims["aud"] = app.APIKey
	t.Claims["jti"] = u.String()
	t.Claims["exp"] = time.Now().Add(5 * time.Minute).Unix()

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
		apiKey := token.Claims["aud"]
		if apiKey == "" {
			return nil, errors.New("api key not found in aud claim")
		}

		//get the public key from the vault
		keystring, err := secretsRepo.RetrievePublicKeyForApp(apiKey.(string))
		if err != nil {
			return nil, err
		}

		//Parse the keystring
		return jwt.ParseRSAPublicKeyFromPEM([]byte(keystring))
	}
}
