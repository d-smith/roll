package secrets
import (
	"crypto/rand"
	"encoding/base64"
	"crypto/rsa"
	"encoding/pem"
	"crypto/x509"
)

func GenerateApiSecret() (string,error) {
	randbuf := make([]byte, 32)

	_,err := rand.Read(randbuf)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(randbuf), nil
}


func GenerateKeyPair()(string, string, error) {
	privateKey ,err  := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return "", "", nil
	}

	privatePEM := pem.EncodeToMemory(
		&pem.Block{
			Type: "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)


	pubkey, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", nil
	}

	publicPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkey,
		},
	)

	return string(privatePEM), string(publicPEM), nil
}

