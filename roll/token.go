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
