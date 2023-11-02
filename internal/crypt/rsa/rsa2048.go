package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

type RSA2048 struct{}

func GenerateKeys() (*Keys, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	return &Keys{key}, nil
}

func (r *RSA2048) Encrypt(plaintext string, key any) (string, error) {
	keys, ok := key.(*Keys)
	if !ok {
		return "", errors.New("invalid encryption key")
	}

	cipherText, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &keys.Private.PublicKey, []byte(plaintext), nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (r *RSA2048) Decrypt(ciphertext string, key any) (string, error) {
	keys, ok := key.(*Keys)
	if !ok {
		return "", errors.New("invalid decryption key")
	}

	rawCipherText, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, keys.Private, rawCipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
