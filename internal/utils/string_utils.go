package utils

import (
	"github.com/jxskiss/base62"
	"math/rand"
)

func RandStr(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func Base62Encode(data string) string {
	return base62.EncodeToString([]byte(data))
}

func Base62Decode(encoded string) (string, error) {
	decodedBytes, err := base62.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}
