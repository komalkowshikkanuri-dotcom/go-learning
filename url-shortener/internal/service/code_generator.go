package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateShortCode(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be positive")
	}

	code := make([]byte, length)

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", fmt.Errorf("generate random index: %w", err)
		}

		code[i] = alphabet[n.Int64()]
	}

	return string(code), nil
}