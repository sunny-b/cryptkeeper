package ecc

import (
	"crypto/elliptic"
	"crypto/sha256"
	"errors"
	"io"

	"github.com/sunny-b/cryptkeeper/internal/crypt/aes"

	"golang.org/x/crypto/hkdf"
)

type ECC256 struct{}

func (e *ECC256) Encrypt(plaintext string, key interface{}) (string, error) {
	k, ok := key.(*Key)
	if !ok {
		return "", errors.New("invalid ecc encryption key")
	}

	// Derive shared secret
	x, _ := elliptic.P256().ScalarMult(k.Private.X, k.Private.Y, k.Private.D.Bytes())
	sharedSecret := x.Bytes()

	// Derive AES key using HKDF
	hkdf := hkdf.New(sha256.New, sharedSecret, nil, nil)
	aesKey := make([]byte, 32)
	_, err := io.ReadFull(hkdf, aesKey)
	if err != nil {
		return "", err
	}

	a := new(aes.AES256)

	return a.Encrypt(plaintext, &aes.EncryptionKey{Key: aesKey})
}

func (e *ECC256) Decrypt(ciphertext string, key interface{}) (string, error) {
	k, ok := key.(*Key)
	if !ok {
		return "", errors.New("invalid decryption key")
	}

	x, _ := elliptic.P256().ScalarMult(k.Private.X, k.Private.Y, k.Private.D.Bytes())
	sharedSecret := x.Bytes()

	// Derive AES key using HKDF
	hkdf := hkdf.New(sha256.New, sharedSecret, nil, nil)
	aesKey := make([]byte, 32)
	_, err := io.ReadFull(hkdf, aesKey)
	if err != nil {
		return "", err
	}

	return new(aes.AES256).Decrypt(ciphertext, &aes.EncryptionKey{Key: aesKey})
}
