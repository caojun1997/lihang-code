package helper

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GetRandomCode(n int) string {
	var sb strings.Builder
	sb.Grow(n)
	for i := 0; i < n; i++ {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		sb.WriteByte(letterBytes[idx.Int64()])
	}
	return sb.String()
}

func RandomString(n int) string {
	return GetRandomCode(n)
}
