package key

import (
	"crypto/rand"
	"crypto/rsa"
)

// GenerateRSAPrivateKey generates an returns an rsa.PrivateKey
func GenerateRSAPrivateKey() (*rsa.PrivateKey, error) {
	pvtKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, NewErrGenerationError(err)
	}

	return pvtKey, nil
}
