package roll

import (
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
