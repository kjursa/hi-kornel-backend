package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"strconv"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const userIdLength = 14
const generalIdLength = 18

func GenerateUserID() string {
	return GenerateSecret(userIdLength)
}
func GenerateID() string {
	return GenerateSecret(generalIdLength)
}

func GenerateSecret(length int) string {
	id := make([]byte, length)

	for i := range length {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return ""
		}
		id[i] = charset[num.Int64()]
	}

	return string(id)
}

func HmacSha256(input string, secret string) string {
	key := []byte(secret)
	message := []byte(input)

	hash := hmac.New(sha256.New, key)
	hash.Write(message)

	return hex.EncodeToString(hash.Sum(nil))
}

func EncodeToBase64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func Int64ToString(n int64) string {
	return strconv.FormatInt(n, 10)
}

func StringToInt64(value string) int64 {
	v, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		panic("Can not parse value: " + value)
	}
	return v
}
