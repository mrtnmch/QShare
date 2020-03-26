package random

import (
	"math/rand"
	"time"
)

var Alphabet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890 !#$%&'()*+,-./:;<=>?@[]^_`{|}~")

func NewStringGenerator() func(n int) string {
	return NewStringGeneratorWithAlphabet(Alphabet)
}

func NewStringGeneratorWithAlphabet(alphabet []byte) func(n int) string {
	rand.Seed(time.Now().UnixNano())

	return func(n int) string {
		return string(generateRandom(alphabet, n))
	}
}

func NewGenerator() func(n int) []byte {
	return NewGeneratorWithAlphabet(Alphabet)
}

func NewGeneratorWithAlphabet(alphabet []byte) func(n int) []byte {
	rand.Seed(time.Now().UnixNano())

	return func(n int) []byte {
		return generateRandom(alphabet, n)
	}
}

func generateRandom(alphabet []byte, n int) []byte {
	if n < 1 {
		return []byte{}
	}

	b := make([]byte, n)

	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return b
}
