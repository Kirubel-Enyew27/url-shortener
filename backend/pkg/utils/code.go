package utils

import "crypto/rand"

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateCode(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "default1"
	}

	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}

	return string(b)
}
