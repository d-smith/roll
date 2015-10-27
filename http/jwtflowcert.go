package http

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/xtraclabs/roll/roll"
	"log"
	"net/http"
	"strings"
)

const (
	//JWTFlowCertsURI is the base uri for the service.
	JWTFlowCertsURI = "/v1/jwtflowcerts/"
)

var (
	errApplicationNotFound = errors.New("No application found for application key")

	errReadingApplicationRecord = errors.New("Error reading application data for application key")

	errInvalidClientSecret = errors.New("")
)

type certPutCtx struct {
	ClientSecret string
	CertPEM string
}

type publicKeyCtx struct {
	PublicKey string
}

func handleJWTFlowCerts(core *roll.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			handleCertPut(core, w, r)
		case "GET":
			handleGetPublicKey(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, errors.New("Method not allowed"))
		}

	})

}

func validateClientSecret(core *roll.Core, r *http.Request, clientID, clientSecret string) (*roll.Application, error) {

	app, err := core.RetrieveApplication(clientID)
	if err != nil {
		return nil, errReadingApplicationRecord
	}

	if app == nil {
		return nil, errApplicationNotFound
	}

	if clientSecret != app.ClientSecret {
		return nil, errInvalidClientSecret
	}

	return app, nil
}

func extractPublicKeyFromCert(certPEM string) (string, error) {
	log.Println("extract public key from:")
	log.Println(certPEM)
	log.Println("certPEM len: ", len(certPEM))

	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return "", errors.New("Unable to decode certificate PEM")
	}

	log.Println("parse the cert")
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", errors.New("failed to parse certificate: " + err.Error())
	}

	pk := cert.PublicKey

	pkbytes, err := x509.MarshalPKIXPublicKey(pk)
	if err != nil {
		return "", errors.New("unable to marshal public key")
	}

	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pkbytes,
		},
	)

	return string(pemdata), nil
}

func checkBodyContent(certCtx certPutCtx) error {
	if certCtx.ClientSecret == "" {
		return errors.New("Request has empty ClientSecret")
	}

	if certCtx.CertPEM == "" {
		return errors.New("Request has empty CertPEM")
	}

	return nil
}

func handleCertPut(core *roll.Core, w http.ResponseWriter, r *http.Request) {
	//Extract client id
	clientID := strings.TrimPrefix(r.RequestURI, JWTFlowCertsURI)
	if clientID == "" {
		respondError(w, http.StatusNotFound, errors.New("Resource not specified"))
		return
	}

	log.Println("Putting cert for client_id", clientID)

	//Parse body
	var certCtx certPutCtx
	if err := parseRequest(r, &certCtx); err != nil {
		log.Println("Error parsing request body", err.Error())
		respondError(w, http.StatusBadRequest, err)
		return
	}

	//Check body content
	err := checkBodyContent(certCtx); if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	//Validate client secret
	app, err := validateClientSecret(core, r, clientID, certCtx.ClientSecret)
	if err != nil {
		switch err {
		case errApplicationNotFound:
			respondNotFound(w)
		case errInvalidClientSecret:
			respondError(w, http.StatusUnauthorized, nil)
		default:
			respondError(w, http.StatusInternalServerError, err)
		}
		return
	}

	//Extract public key from cert
	publicKeyPEM, err := extractPublicKeyFromCert(certCtx.CertPEM)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	//Update the app with the public key. Note here we are adding the cert to the retrieved application
	//attributes.
	app.JWTFlowPublicKey = publicKeyPEM
	err = core.UpdateApplication(app)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, nil)
}

func handleGetPublicKey(core *roll.Core, w http.ResponseWriter, r *http.Request) {
	//Extract client id
	clientID := strings.TrimPrefix(r.RequestURI, JWTFlowCertsURI)
	if clientID == "" {
		respondError(w, http.StatusNotFound, errors.New("Resource not specified"))
		return
	}

	//Retrieve the app definition
	app, err := core.RetrieveApplication(clientID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, errReadingApplicationRecord)
		return
	}

	if app == nil {
		respondError(w, http.StatusNotFound, nil)
		return
	}

	pk := publicKeyCtx{
		PublicKey:app.JWTFlowPublicKey,
	}

	respondOk(w,&pk)
}
