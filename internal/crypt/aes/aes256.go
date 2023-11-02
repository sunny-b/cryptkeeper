package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

type EncryptionKey struct {
	Key []byte `json:"key"`
}

type AES256 struct {
	aead cipher.AEAD
}

func GenerateKeys() (*EncryptionKey, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}

	return &EncryptionKey{key}, nil
}

// Encrypt encrypts the given plaintext with the given key using AES-256-GCM
func (a *AES256) Encrypt(plaintext string, key any) (string, error) {
	e, ok := key.(*EncryptionKey)
	if !ok {
		return "", errors.New("invalid encryption key")
	}

	if a.aead == nil {
		block, err := aes.NewCipher(e.Key)
		if err != nil {
			return "", err
		}

		a.aead, err = cipher.NewGCM(block)
		if err != nil {
			return "", err
		}
	}

	nonce := make([]byte, a.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := a.aead.Seal(nonce, nonce, []byte(plaintext), nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts the given ciphertext with the given key using AES-256-GCM
func (a *AES256) Decrypt(cipherText string, key any) (string, error) {
	e, ok := key.(*EncryptionKey)
	if !ok {
		return "", errors.New("invalid encryption key")
	}

	rawCiphertext, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	if a.aead == nil {
		block, err := aes.NewCipher(e.Key)
		if err != nil {
			return "", err
		}

		a.aead, err = cipher.NewGCM(block)
		if err != nil {
			return "", err
		}
	}

	nonceSize := a.aead.NonceSize()
	if len(rawCiphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, cipher := rawCiphertext[:nonceSize], rawCiphertext[nonceSize:]

	b, err := a.aead.Open(nil, nonce, cipher, nil)

	return string(b), err
}
