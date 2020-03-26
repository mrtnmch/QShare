package qr

import (
	"fmt"
	"github.com/skip2/go-qrcode"
)

type Image struct {
	Content []byte
}

type Generator interface {
	NewPNGImage(content []byte) (*Image, error)
}

type DefaultGenerator struct {
}

func NewGenerator() *DefaultGenerator {
	return &DefaultGenerator{}
}

func (generator *DefaultGenerator) NewPNGImage(content []byte) (*Image, error) {
	png, err := qrcode.Encode(string(content), qrcode.Highest, len(string(content)))

	if err != nil {
		return nil, fmt.Errorf("generate QR file %v: %v", content, err)
	}

	image := &Image{
		Content: png,
	}

	return image, nil
}