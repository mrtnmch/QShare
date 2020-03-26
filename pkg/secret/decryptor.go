package secret

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

type Decryptor interface {
	Decrypt(content []byte) ([]byte, error)
}

type RSADecryptor struct {
	privateKey *rsa.PrivateKey
}

func NewRSADecryptor(privateKey *rsa.PrivateKey) *RSADecryptor {
	return &RSADecryptor{
		privateKey,
	}
}

func (decryptor RSADecryptor) Decrypt(content []byte) ([]byte, error) {
	return rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		decryptor.privateKey,
		content,
		[]byte(""))
}
