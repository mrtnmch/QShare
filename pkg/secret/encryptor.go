package secret

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

type Encryptor interface {
	Encrypt(content []byte) ([]byte, error)
}

type RSAEncryptor struct {
	publicKey  *rsa.PublicKey
}

func NewRSAEncryptor(publicKey *rsa.PublicKey, secret []byte) *RSAEncryptor {
	return &RSAEncryptor{
		publicKey:  publicKey,
	}
}

func (encryptor *RSAEncryptor) Encrypt(content []byte) ([]byte, error) {
	encrypted, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		encryptor.publicKey,
		content,
		[]byte(""))

	if err != nil {
		return nil, fmt.Errorf("error has occurred while encrypting the message: %v", err)
	}

	return encrypted, nil
}