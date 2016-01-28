package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	rollhttp "github.com/xtraclabs/roll/http"
	"log"
	"net/http"
)

const auth0Cert = `
-----BEGIN CERTIFICATE-----
MIIDBTCCAe2gAwIBAgIJAKLDWbFECKaKMA0GCSqGSIb3DQEBCwUAMBkxFzAVBgNV
BAMMDnhhdmkuYXV0aDAuY29tMB4XDTE1MDgxMjE5MzUwMFoXDTI5MDQyMDE5MzUw
MFowGTEXMBUGA1UEAwwOeGF2aS5hdXRoMC5jb20wggEiMA0GCSqGSIb3DQEBAQUA
A4IBDwAwggEKAoIBAQDcoYmMt282UITLAz240hjvQiGap0JIkwUMKgRiaQ5VLwWv
YAVORM2fhVCuZp4gSGCyF/FFuksXQ7ONiv5CysvPH+Msy0xH6A2ugmcp3LYGPi5O
BRajDdC+Evbm8xoPJIJcxsEoVbIJlf0P3dihL78H412n18oKlieQg9zY5t48BD1C
s/qaIu9tN7SvYCPYnkM5jBOsj9yxIpyPLsFeRJ4gVefOoOFTP/Uramc9y0wUwz+p
wXgBFoR3tZABx3nsLxVHtf29KeWmSz9ogu6N41Sw5dOUpMFptfYZVSlyz74G3myN
8t4Zk/cddMyLPeZzKmdkJmR7j7VoIVSkM6L8JmN7AgMBAAGjUDBOMB0GA1UdDgQW
BBTyET2FNhUsEPN6uSK8Hrp8qPdmyjAfBgNVHSMEGDAWgBTyET2FNhUsEPN6uSK8
Hrp8qPdmyjAMBgNVHRMEBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQCddBfJ5tVJ
mdfljSrJPj2kZn/o+ECi0Lzap+rtDypuJrcHyFJzxZOAryDTrRDgG4Xo9ZJgtCYI
9UQxJ+fa3DRAnqMJ0lZUrS5Vmlkmt3CWI8BZyc8fcgJj3vIyUn6qyS6+Z3FhnlPL
trdwUKjIcFoctCOMQIYCGYXW3qu5YQ+xc0mdopYs7lsPWdG9D9fdhZmSagOzx9pr
BdQugpK2yy1nfX9f66gg+NbNFkNjPX8ff0Q8uzBrUlG0LqgVEs2g/VHZgNq5PQLA
R5MNabBfadFN1fa8S4+5MlUoZJGaIF8YVVgM8pV880t8zcP3MZ3cm83z8yNw2Exl
xYOAhtYv9gls
-----END CERTIFICATE-----

`

const (
	clientID       = "d80c9758-f965-4286-7003-e5af065d8082"
	clientSecret   = "4suKio+75Vsm2o+i8H8dEjxVbXhoB81G+rJolRw8+XU="
	baseJWTCertURL = "http://localhost:3000/v1/jwtflowcerts/"
	tokenURL       = "http://localhost:3000/oauth2/token"
	authToken      = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBsaWNhdGlvbiI6InNhbXBsZSBkZXYgYXBwIiwiYXVkIjoiZDgwYzk3NTgtZjk2NS00Mjg2LTcwMDMtZTVhZjA2NWQ4MDgyIiwiZXhwIjoxNDU0MDc3NjE0LCJpYXQiOjE0NTM5OTEyMTQsImp0aSI6IjM5YzIzZTZiLWRjYmEtNDVmOC01OTA4LTcyYmEzZTZlZTllNyIsInNjb3BlIjoiIiwic3ViIjoidXNlciJ9.AelDKwAoJny_hbgLzeRZYbfLyNVi9A2csMIp8pYirp6EbSp-n7akvgmKKV3Rn_ZiElZRbkg4safXVRVXCT-pirHAn_7pzwI2HthNaqWIKLdjASo-PkZmlEgv_9Sbprpn9JTsQfgVUXNjnaKqNWSWYd3qCb43TrK-9TMYewKYaVQ"
)

func uploadAuth0Cert(theCert string) {
	fmt.Println("Uploading cert to", baseJWTCertURL+clientID)
	fmt.Println(theCert)

	payload := rollhttp.CertPutCtx{
		ClientSecret: clientSecret,
		CertPEM:      theCert,
		CertIssuer:   "https://xavi.auth0.com/",
		CertAudience: "vY0bFoxCBzE9rrTNTEjhIfay8MbFYq9Z",
	}

	bodyReader := new(bytes.Buffer)
	enc := json.NewEncoder(bodyReader)
	err := enc.Encode(&payload)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("PUT", baseJWTCertURL+clientID, bodyReader)
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

func main() {
	uploadAuth0Cert(auth0Cert)
}
