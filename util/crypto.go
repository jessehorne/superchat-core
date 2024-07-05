package util

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
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

// GenerateHsh takes a password and salt byte array and returns a byte array containing an argon2 key
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

// CreateToken creates a session-based auth token and sha256 hash of it and returns them
// The token and token hash are base64 encoded.
func CreateToken() (string, string) {
	b := make([]byte, 32)
	rand.Read(b)
	h := sha256.New()
	h.Write(b)
	bs := h.Sum(nil)

	tokenBase64 := base64.RawStdEncoding.EncodeToString(b)
	hashBase64 := base64.RawStdEncoding.EncodeToString(bs)

	return tokenBase64, hashBase64
}

// ValidateToken takes base64 encoded token and token hash and returns true if they are equal when the token is decoded
// and sha256'd.
func ValidateToken(token string, hash string) bool {
	decodedToken, _ := base64.RawStdEncoding.DecodeString(token)
	decodedHash, _ := base64.RawStdEncoding.DecodeString(hash)

	h := sha256.New()
	h.Write(decodedToken)
	bs := h.Sum(nil)

	return bytes.Equal(bs, decodedHash)
}
