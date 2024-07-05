package util

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"golang.org/x/crypto/argon2"
)

const (
	defaultSaltLength = 16
)

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var argonParams = &params{
	memory:      64 * 1024,
	iterations:  3,
	parallelism: 2,
	saltLength:  defaultSaltLength,
	keyLength:   32,
}

func GenerateSalt(len int) []byte {
	salt := make([]byte, len)
	rand.Read(salt)
	return salt
}

func GenerateHash(pass string, salt []byte) []byte {
	return argon2.IDKey([]byte(pass), salt,
		argonParams.iterations, argonParams.memory, argonParams.parallelism, argonParams.keyLength)
}

// ProcessPassword returns the base64 encoded salt and base64 encoded argon2 hash or an error
func ProcessPassword(pass string) (string, string) {
	salt := GenerateSalt(defaultSaltLength)
	saltString := base64.RawStdEncoding.EncodeToString(salt)

	hash := GenerateHash(pass, salt)
	hashString := base64.RawStdEncoding.EncodeToString(hash)

	return saltString, hashString
}

// ComparePassword compares pass with salt and hash and returns true if it's valid
// You must provide the base64 encoded salt and hash to this function.
func ComparePassword(pass string, salt string, hash string) bool {
	decodedSalt, _ := base64.RawStdEncoding.DecodeString(salt)
	decodedHash, _ := base64.RawStdEncoding.DecodeString(hash)

	generatedHash := GenerateHash(pass, decodedSalt)

	if subtle.ConstantTimeCompare(decodedHash, generatedHash) == 1 {
		return true
	}

	return false
}
