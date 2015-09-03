package roll

import (
	jwt "github.com/dgrijalva/jwt-go"
)


func GenerateToken(app *Application, privateKey string) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	t.Claims["APIKey"] = app.APIKey
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return "",nil
	}

	tokenString,err := t.SignedString(signKey)
	if err != nil {
		return "",err
	}

	return tokenString, nil
}