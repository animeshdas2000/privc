package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
)

// Token layout (before base64url encoding, no padding):
//
//	version (1 byte) | nonce (12 bytes) | AES-GCM ciphertext + 16-byte tag
//
// The version byte is also bound as additional authenticated data, so it
// cannot be swapped without failing authentication.
const tokenVersion byte = 1

var errInvalidToken = errors.New("invalid token")

func newGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

// Encrypt seals plaintext with AES-GCM under a fresh random nonce, so the
// same input produces a different token every time.
func Encrypt(key []byte, plaintext string) (string, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	out := make([]byte, 0, 1+len(nonce)+len(plaintext)+gcm.Overhead())
	out = append(out, tokenVersion)
	out = append(out, nonce...)
	out = gcm.Seal(out, nonce, []byte(plaintext), []byte{tokenVersion})
	return base64.RawURLEncoding.EncodeToString(out), nil
}

// Decrypt opens a token produced by Encrypt. Any malformed, tampered or
// foreign token returns errInvalidToken.
func Decrypt(key []byte, token string) (string, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return "", err
	}
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return "", errInvalidToken
	}
	ns := gcm.NonceSize()
	if len(raw) < 1+ns+gcm.Overhead() || raw[0] != tokenVersion {
		return "", errInvalidToken
	}
	nonce, sealed := raw[1:1+ns], raw[1+ns:]
	plaintext, err := gcm.Open(nil, nonce, sealed, raw[:1])
	if err != nil {
		return "", errInvalidToken
	}
	return string(plaintext), nil
}
