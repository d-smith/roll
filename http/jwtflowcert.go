package http

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/xtraclabs/roll/roll"
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

type CertPutCtx struct {
	ClientSecret string `json:"clientSecret"`
	CertPEM      string `json:"certPEM"`
	CertIssuer   string `json:"issuer`
	CertAudience string `json:"audience"`
}

type publicKeyCtx struct {
	PublicKey string `json:"publicKey"`
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

	app, err := core.SystemRetrieveApplication(clientID)
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
	log.Info("extract public key from:")
	log.Info(certPEM)
	log.Info("certPEM len: ", len(certPEM))

	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return "", errors.New("Unable to decode certificate PEM")
	}

	log.Info("parse the cert")
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

func checkBodyContent(certCtx CertPutCtx) error {
	if certCtx.ClientSecret == "" {
		return errors.New("Request has empty ClientSecret")
	}

	if certCtx.CertPEM == "" {
		return errors.New("Request has empty CertPEM")
	}

	if certCtx.CertIssuer == "" {
		return errors.New("Request has empty Issuer")
	}

	if certCtx.CertAudience == "" {
		return errors.New("Request has empty audience")
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

	log.Info("Putting cert for client_id: ", clientID)

	//Extract the subject from the request header based on security mode
	subject, _, err := subjectAndAdminScopeFromRequestCtx(r)
	if err != nil {
		log.Print("Error extracting subject: ", err.Error())
		respondError(w, http.StatusInternalServerError, nil)
		return
	}

	//Parse body
	var certCtx CertPutCtx
	if err := parseRequest(r, &certCtx); err != nil {
		log.Info("Error parsing request body: ", err.Error())
		respondError(w, http.StatusBadRequest, err)
		return
	}

	//Check body content
	log.Info("Checking content")
	err = checkBodyContent(certCtx)
	if err != nil {
		log.Info("Problem with content: ", err.Error())
		respondError(w, http.StatusBadRequest, err)
		return
	}

	//Validate client secret
	log.Info("validating client secret")
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
	log.Info("Extract public key")
	publicKeyPEM, err := extractPublicKeyFromCert(certCtx.CertPEM)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	//Update the app with the public key. Note here we are adding the cert to the retrieved application
	//attributes.
	log.Info("Update app with signing key, etc")
	app.JWTFlowPublicKey = publicKeyPEM
	app.JWTFlowIssuer = certCtx.CertIssuer
	app.JWTFlowAudience = certCtx.CertAudience
	err = core.UpdateApplication(app, subject)
	if err != nil {
		switch err.(type) {
		case roll.NonOwnerUpdateError:
			respondError(w, http.StatusUnauthorized, err)
		case roll.NoSuchApplicationError:
			respondError(w, http.StatusNotFound, err)
		case roll.MissingJWTFlowIssuer:
			respondError(w, http.StatusBadRequest, err)
		case roll.MissingJWTFlowAudience:
			respondError(w, http.StatusBadRequest, err)
		default:
			respondError(w, http.StatusInternalServerError, err)
		}

		return
	}

	respondOk(w, nil)
}

func handleGetPublicKey(core *roll.Core, w http.ResponseWriter, r *http.Request) {
	//Extract client id
	clientID := strings.TrimPrefix(r.RequestURI, JWTFlowCertsURI)
	if clientID == "" {
		respondError(w, http.StatusBadRequest, errors.New("Resource not specified"))
		return
	}

	log.Info("retrieve public key for application: ", clientID)

	//Retrieve the app definition. Note that here since we are only returning publically
	//available information, we do not have to apply the data security model
	app, err := core.SystemRetrieveApplication(clientID)
	if err != nil {
		log.Info("error retrieving application")
		respondError(w, http.StatusInternalServerError, errReadingApplicationRecord)
		return
	}

	if app == nil {
		log.Info("application not found")
		respondError(w, http.StatusNotFound, nil)
		return
	}

	pk := publicKeyCtx{
		PublicKey: app.JWTFlowPublicKey,
	}

	respondOk(w, &pk)
}
