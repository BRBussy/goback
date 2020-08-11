package key

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

// ParseRSAPrivateKeyFromString parses parses an RSA private key in PKCS#1, ASN.1 DER form.
// It returns an a pointer to an rsa.PrivateKey which represents an RSA Private Key.
func ParseRSAPrivateKeyFromString(rsaPrivateKeyString string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(rsaPrivateKeyString))
	if block == nil {
		return nil, NewErrNilKey()
	}

	pvtKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, NewErrParsingError(err)
	}

	return pvtKey, nil
}

// GenerateRSAPrivateKey generates an returns an rsa.PrivateKey
func GenerateRSAPrivateKey() (*rsa.PrivateKey, error) {
	pvtKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, NewErrGenerationError(err)
	}

	return pvtKey, nil
}
