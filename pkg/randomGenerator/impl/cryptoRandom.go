package impl

import (
	"2022_1_OnlyGroup_back/app/handlers/http"
	"crypto/rand"
	"encoding/base32"
)

type cryptoRandomGenerator struct {
}

func NewCryptoRandomGenerator() *cryptoRandomGenerator {
	return &cryptoRandomGenerator{}
}

func (*cryptoRandomGenerator) String(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", http.ErrBaseApp.Wrap(err, "crypto rand failed")
	}

	return base32.StdEncoding.EncodeToString(randomBytes)[:length], nil
}

func (*cryptoRandomGenerator) Bytes(length int) ([]byte, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return []byte{}, http.ErrBaseApp.Wrap(err, "crypto rand failed")
	}
	return randomBytes, nil
}
