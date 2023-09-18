package crypt

import "errors"

type EncryptionType string

var (
	ErrUnknownEncryptionType                = errors.New("unknown encryption type")
	AES256                   EncryptionType = "aes256"
	ECC256                   EncryptionType = "ecc256"
	RSA2048                  EncryptionType = "rsa2048"
	Serpent256               EncryptionType = "serpent256"
)
