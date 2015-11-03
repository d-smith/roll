package main

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

const certPEM = `
-----BEGIN CERTIFICATE-----
MIIC+DCCAeKgAwIBAgIRAIJaB8pAErenO9pMBUDo3awwCwYJKoZIhvcNAQELMBIx
EDAOBgNVBAoTB0FjbWUgQ28wHhcNMTUwODI4MTM0MzU3WhcNMTYwODI3MTM0MzU3
WjASMRAwDgYDVQQKEwdBY21lIENvMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB
CgKCAQEAwoH4nc/B3/i1D1TjCx3kgC6ygX3WHDv/xHtAoRAgHFUVElo3PznbxLAk
MvElVdAevCJaiuJaLiZARKLvwSJh08/9y+WMYa1nDjINk6UqG3huPXdJmTguzleO
c7UrCW4WKSo2HbeqYlF4BOiqnQhdDncUh5BgR8JXuiueMn2Ka59lkB/i+ryOt5W7
kaFKJhQEV67+fuES/5WfE+B4XsfT/ctXnGY0zrEInbJlyKwAzyCWJOJFrZte8cxs
235q3VMAhMRDU1IGNuWBIntfEXZgUXqI1Z9gsdbfTsQQ+xWhQCCOJwDrxAEg1Udk
dWn6NGWevsH4JoM9JzzOeSH8ZYPrVQIDAQABo00wSzAOBgNVHQ8BAf8EBAMCAKAw
EwYDVR0lBAwwCgYIKwYBBQUHAwEwDAYDVR0TAQH/BAIwADAWBgNVHREEDzANggtN
QUNMQjAxNTgwMzALBgkqhkiG9w0BAQsDggEBALJtJGaXx9At98CvEWKBpiGYqjUu
aiQHS5R61R/g8iqWkct77cqN6SBWTf138NZ3j3mvfROCoU96BEMEl0Fk9apLrikI
9Ns9/sl4nL1IOR56vddm46DfEV5CpMCAgrMGhFMJiaW4t9HvYjpBSs8T5n4tGqu/
JsvPhLGOcu5i4RiPpwM8f4fhnD3sija334jj5meJwg0NR8eO3ro1zaH+0hMQ7l8Q
tFJusSJenG28q9MXpOoCG6KLCmSCrIfDRYIpJQ0d5fXLO4YG92KFFqrf2ycOTydY
hN9G5ZWaErEY5j+sbYmeJBtEM5v6BQJotJh2SAh8RpYr69qJPLw6fdTu+mU=
-----END CERTIFICATE-----`

const keyPEM = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAwoH4nc/B3/i1D1TjCx3kgC6ygX3WHDv/xHtAoRAgHFUVElo3
PznbxLAkMvElVdAevCJaiuJaLiZARKLvwSJh08/9y+WMYa1nDjINk6UqG3huPXdJ
mTguzleOc7UrCW4WKSo2HbeqYlF4BOiqnQhdDncUh5BgR8JXuiueMn2Ka59lkB/i
+ryOt5W7kaFKJhQEV67+fuES/5WfE+B4XsfT/ctXnGY0zrEInbJlyKwAzyCWJOJF
rZte8cxs235q3VMAhMRDU1IGNuWBIntfEXZgUXqI1Z9gsdbfTsQQ+xWhQCCOJwDr
xAEg1UdkdWn6NGWevsH4JoM9JzzOeSH8ZYPrVQIDAQABAoIBAQC7HhvVfZNuaq2M
RUZyYFx9L+9MAYmj+dhw5X1kNMYfW358/EsfyQJx8l2woCohpqiappXh/CVzNhsn
rHPnobLo+LOYnojxQsbccix7eC3M5YtwGDy2aMUmfjate4EefkhVa6VJ63G/Sc6j
Mk1H6k7TZp1aWPs2rXCFygZTXlB2fOCpGRVfTzKO8p+x3mpLw54oPzoKDpiu2g0a
gMIXWlqv+nVNx7TT0o0koENnHyIlFE87ILvuH2FV6o6jMC7W84WIv+3LfSjEH3Np
VlYL7e9DCqK/rvSjZo7F5ZV8CUn/zfPOpS2IHfKak1Qcar1o03bX8XUY1HNcYrz3
OUg+30NhAoGBANbHeuG4ZLSSEUmxT/JdhHb6XxvTbImXgyzed6Qvcpz/qY2ycijK
ed+sh7u7QoNj+UCb89pHwg8staBLGZULwTPXR7HIC8UA+u7CI/NAkOLs/T9fBj9E
mB/bHgnLM6cID08Uqhc/Rw+YGH2AFGrCNPI8IY/m3+aK8AoBo/WEXxrdAoGBAOfW
gxcLxr97HHY+SNbxuNXS8jeuuKf3wr9fFSPqarkResCYGjqj/0i25sPc+gh/Jksb
GJAVr5rPYf0NEjVxJffITuHONhrwidEgqYrejUaGCpqmo4H/1MHDrccns6nj955J
fQrd1wDaV/y6Pi9ZTddtfZvgTu7F4wDOLqKCHV7ZAoGABQ9eDcNSXA606Y+L2qHp
ZQQxCW7Jt9Qh+wSivPOBN8GuNPBnUB61EFf4OqwGvSG11ZfW/uWB5Omgvc4HJEPK
oXYI8l9MJ7AGOWnJfwmeYDpFHSEX/JWZH+cBpilrZhwCgFJ0T2fmPK/Qo+RIlttM
dAiGoeVGkR/k19eCDV64/uECgYAEgqTcPutn8EI2frG2FM/OwSLbOHf5NJXCWVw0
7auT3CYyNwNTnrsps/07VEI9BhQ9o9JPg0x+S4iRfr6a8Y+WiXwXYoak8YcM3YDO
5dODyrnMFm/Z6JtABcmMtVYGeRvpJW7cVkKHjrflf84T573e8g9eh/ehJ5JqApjf
qn5dyQKBgFnFGPBaglMn0tMZ1E61FKFpa5dOAybvDsUIom58oJKy6lBMouKpI/sK
8qALTOLVRz36zsNCkjfx9SlXCvZV2CA5MXH/jKWOGbZrhOpOHHh/lCPb9+bISKoW
T9276oM42khyKY36lXvLi4yjk2yHysIvO7ckuX0F/vZtQjG1zuBb
-----END RSA PRIVATE KEY-----
`

const (
	clientID     = "b047f0b9-e830-4c6b-7c0c-0564bd07a0a7"
	clientSecret = "ThWLhQYfW1s/9aXEA5W5Lszx1UoZUWebVv4j1KWPKxk="
)

func uploadCert() {
	fmt.Println("Uploading cert")
	fmt.Println(certPEM)

	payload := rollhttp.CertPutCtx{
		ClientSecret: clientSecret,
		CertPEM:      certPEM,
	}

	bodyReader := new(bytes.Buffer)
	enc := json.NewEncoder(bodyReader)
	err := enc.Encode(&payload)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "http://localhost:3000/v1/jwtflowcerts/"+clientID, bodyReader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusNoContent {
		log.Fatal("Did not receive an OK for cert upload, got ", resp.StatusCode)
	}
}

func generateJTW() string {
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(keyPEM))
	if err != nil {
		log.Fatal("Unable to parse the key PEM")
	}

	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims["iss"] = clientID
	token.Claims["sub"] = "drscan"

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		log.Fatal("Unable to sign token: ", err.Error())
	}

	return tokenString

}

func tradeTokenForToken(token string) string {
	resp, err := http.PostForm("http://localhost:3000/oauth2/token",
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

func main() {
	uploadCert()
	tokenString := generateJTW()
	fmt.Println("\nUse ", tokenString, " to obtain access token")
	jwtResponse := tradeTokenForToken(tokenString)
	fmt.Println("\n", jwtResponse)
}
