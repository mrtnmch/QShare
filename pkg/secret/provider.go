package secret

import "encoding/base64"

type Provider interface {
	Generate() ([]byte, error)
}

type StaticProvider struct {
	encryptor Encryptor
	secret []byte
}

func NewStaticProvider(encryptor Encryptor, secret []byte) *StaticProvider {
	return &StaticProvider{
		encryptor: encryptor,
		secret:    secret,
	}
}

func (provider *StaticProvider) Generate() ([]byte, error) {
	return provider.encryptor.Encrypt(provider.secret)
}

type Base64Provider struct {
	upstream Provider
}

func NewBase64Provider(provider Provider) *Base64Provider {
	return &Base64Provider{
		upstream:provider,
	}
}

func (provider *Base64Provider) Generate() ([]byte, error) {
	generate, err := provider.upstream.Generate()
	target := base64.StdEncoding.EncodeToString(generate)
	return []byte(target), err
}