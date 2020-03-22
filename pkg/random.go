package pkg

import (
	"math/rand"
	"time"
)

var Alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func GenerateRandomString(n int) string {
	if n < 1 {
		return ""
	}

	b := make([]rune, n)
	rand.Seed(time.Now().UnixNano())

	for i := range b {
		b[i] = Alphabet[rand.Intn(len(Alphabet))]
	}

	return string(b)
}
