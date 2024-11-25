package url_generator

import (
	"math/rand"
)

const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateShortURL(length int) string {
	shortURL := make([]byte, length)
	for i := range shortURL {
		shortURL[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(shortURL)
}
