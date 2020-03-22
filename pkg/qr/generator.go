package qr

import (
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/skip2/go-qrcode"
)

type Generator interface {
	Generate() (*Image, error)
	Validate(base64code string) error
}

type Image struct {
	Code        string
	Base64Image string
}

type DefaultGenerator struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	secret     string
}

func NewGenerator(secret string, key *rsa.PrivateKey) *DefaultGenerator {
	generator := &DefaultGenerator{
		publicKey:  &key.PublicKey,
		privateKey: key,
		secret:     secret,
	}

	return generator
}

func (generator *DefaultGenerator) Generate() (*Image, error) {
	var png []byte

	encrypted, err := rsa.EncryptOAEP(
		sha256.New(),
		cryptorand.Reader,
		generator.publicKey,
		[]byte(generator.secret),
		[]byte(""))

	if err != nil {
		return nil, errors.New("error has occurred while encrypting the message: " + err.Error())
	}

	code := base64.StdEncoding.EncodeToString(encrypted)
	png, _ = qrcode.Encode(code, qrcode.Highest, len(code))

	image := &Image{
		Code:        code,
		Base64Image: base64.StdEncoding.EncodeToString(png),
	}

	return image, nil
}

func (generator *DefaultGenerator) Validate(base64code string) error {
	hash := sha256.New()
	label := []byte("")

	code, _ := base64.StdEncoding.DecodeString(base64code)

	bytes, err := rsa.DecryptOAEP(
		hash,
		cryptorand.Reader,
		generator.privateKey,
		code,
		label)

	if err != nil {
		return err
	}

	if string(bytes) != generator.secret {
		return errors.New("control secrets don't match")
	}

	return nil
}
