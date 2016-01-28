package jwtsample

import (
	"bytes"
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	rollhttp "github.com/xtraclabs/roll/http"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type RollContext struct {
	BaseJWTCertURL string
	ClientID       string
	ClientSecret   string
	CertPEM        string
}

func UploadCert(rollCtx RollContext, authToken string) {
	fmt.Println("Uploading cert to", rollCtx.BaseJWTCertURL+rollCtx.ClientID)
	fmt.Println(rollCtx)

	payload := rollhttp.CertPutCtx{
		ClientSecret: rollCtx.ClientSecret,
		CertPEM:      rollCtx.CertPEM,
	}

	bodyReader := new(bytes.Buffer)
	enc := json.NewEncoder(bodyReader)
	err := enc.Encode(&payload)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("PUT", rollCtx.BaseJWTCertURL+rollCtx.ClientID, bodyReader)
	req.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusNoContent {
		log.Fatal("Did not receive an OK for cert upload, got ", resp.StatusCode)
	}
}

func GenerateJWT(keyPEM string, clientID string) string {
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(keyPEM))
	if err != nil {
		log.Fatal("Unable to parse the key PEM")
	}

	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims["iss"] = clientID
	token.Claims["sub"] = "foo"
	token.Claims["scope"] = "admin"

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		log.Fatal("Unable to sign token: ", err.Error())
	}

	return tokenString

}

func TradeTokenForToken(token string, tokenURL string) string {
	resp, err := http.PostForm(tokenURL,
		url.Values{"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"},
			"assertion": {token}})
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}
