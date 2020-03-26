package secret

import (
	"bytes"
	"encoding/base64"
	"log"
)

type Validator interface {
	IsValid(content []byte) bool
}

type DefaultValidator struct {
	Decryptor Decryptor
	secret    []byte
}

func NewBase64SecretValidator(secret []byte, decryptor Decryptor) *DefaultValidator {
	return &DefaultValidator{
		Decryptor: decryptor,
		secret:    secret,
	}
}

func (validator *DefaultValidator) IsValid(content []byte) bool {
	decoded, err := base64.StdEncoding.DecodeString(string(content))

	if err != nil {
		log.Println(err)
		return false
	}

	decrypted, err := validator.Decryptor.Decrypt(decoded)

	if err != nil {
		log.Println(err)
		return false
	}

	return bytes.Equal(validator.secret, decrypted)
}
