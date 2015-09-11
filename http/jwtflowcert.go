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
	//JWTCertsURI is the base uri for the service.
	JWTFlowCertsURI = "/v1/jwtflowcerts/"
)

var (
	ErrApplicationNotFound = errors.New("No application found for application key")

	ErrReadingApplicationRecord = errors.New("Error reading application data for application key")

	ErrInvalidClientSecret = errors.New("")
)

type certPostCtx struct {
	clientSecret string
	certPEM      string
}

func handleJWTFlowCerts(core *roll.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			handleCertPost(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, errors.New("Method not allowed"))
		}

	})

}

func extractFormParams(r *http.Request) (*certPostCtx, error) {
	clientSecret := r.FormValue("client_secret")
	if clientSecret == "" {
		return nil, errors.New("client_secret missing from request")
	}

	certPEM := r.FormValue("cert_pem")
	if certPEM == "" {
		return nil, errors.New("cert_pem missing from request")
	}

	return &certPostCtx{
		clientSecret: clientSecret,
		certPEM:      certPEM,
	}, nil
}

func validateClientSecret(core *roll.Core, r *http.Request, clientSecret string) (*roll.Application, error) {
	apiKey := strings.TrimPrefix(r.RequestURI, JWTFlowCertsURI)
	if apiKey == "" {
		return nil, ErrApplicationNotFound
	}

	app, err := core.RetrieveApplication(apiKey)
	if err != nil {
		return nil, ErrReadingApplicationRecord
	}

	if app == nil {
		return nil, ErrApplicationNotFound
	}

	if clientSecret != app.APISecret {
		return nil, ErrInvalidClientSecret
	}

	return app, nil
}

func extractPublicKeyFromCert(certPEM string) (string, error) {
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

func handleCertPost(core *roll.Core, w http.ResponseWriter, r *http.Request) {
	//Extract form parameters
	certCtx, err := extractFormParams(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	//Validate client secret
	app, err := validateClientSecret(core, r, certCtx.clientSecret)
	if err != nil {
		switch err {
		case ErrApplicationNotFound:
			respondNotFound(w)
		case ErrInvalidClientSecret:
			respondError(w, http.StatusUnauthorized, nil)
		default:
			respondError(w, http.StatusInternalServerError, err)
		}
		return
	}

	//Extract public key from cert
	publicKeyPEM, err := extractPublicKeyFromCert(certCtx.certPEM)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	//Update the app with the public key
	app.JWTFlowPublicKey = publicKeyPEM
	err = core.StoreApplication(app)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, nil)
}
