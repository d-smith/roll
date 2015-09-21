package secrets

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

//GenerateClientSecret generates a string to be used as an API secret
func GenerateClientSecret() (string, error) {
	randbuf := make([]byte, 32)

	_, err := rand.Read(randbuf)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(randbuf), nil
}

//GenerateKeyPair generates a random 1024 byte RSA private and public key pair, returning
//PEM encodings of the keys
func GenerateKeyPair() (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return "", "", nil
	}

	privatePEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
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
