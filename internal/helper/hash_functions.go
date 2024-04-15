package helper

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
)

func Sha512HmacHash(input []byte, secretKey string) string {
	// Convert the input string and secret key to bytes
	data := input

	// Create a new HMAC-SHA-512 hash instance.
	hash := hmac.New(sha512.New, []byte(secretKey))

	// Write the input data to the hash.
	hash.Write(data)

	// Calculate the HMAC-SHA-512 hash sum.
	hashSum := hash.Sum(nil)

	// Convert the hash sum to a hexadecimal string.
	hashString := hex.EncodeToString(hashSum)

	return hashString
}
