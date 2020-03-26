package qr

import "encoding/base64"

type Base64Generator struct {
	BaseGenerator Generator
}

func (generator *Base64Generator) NewPNGImage(content []byte) (*Image, error) {
	var encoded []byte
	base64.StdEncoding.Encode(encoded, content)
	return generator.BaseGenerator.NewPNGImage(encoded)
}

func NewBase64Generator(generator Generator) *Base64Generator {
	if generator == nil {
		return nil
	}

	return &Base64Generator{generator}
}


