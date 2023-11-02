package serpent

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/aead/serpent"
)

type EncryptionKey struct {
	Key []byte `json:"key"`
}

type Serpent256 struct {
	aead cipher.AEAD
}

func GenerateKeys() (*EncryptionKey, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}

	return &EncryptionKey{key}, nil
}

func (s *Serpent256) Encrypt(plainText string, key any) (string, error) {
	e, ok := key.(*EncryptionKey)
	if !ok {
		return "", errors.New("invalid encryption key")
	}

	if s.aead == nil {
		block, err := serpent.NewCipher(e.Key)
		if err != nil {
			return "", err
		}

		s.aead, err = cipher.NewGCM(block)
		if err != nil {
			return "", err
		}
	}

	nonce := make([]byte, s.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := s.aead.Seal(nonce, nonce, []byte(plainText), nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *Serpent256) Decrypt(cipherText string, key any) (string, error) {
	e, ok := key.(*EncryptionKey)
	if !ok {
		return "", errors.New("invalid encryption key")
	}

	rawCiphertext, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	if s.aead == nil {
		block, err := serpent.NewCipher(e.Key)
		if err != nil {
			return "", err
		}

		s.aead, err = cipher.NewGCM(block)
		if err != nil {
			return "", err
		}
	}

	nonceSize := s.aead.NonceSize()
	if len(rawCiphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, cipher := rawCiphertext[:nonceSize], rawCiphertext[nonceSize:]

	b, err := s.aead.Open(nil, nonce, cipher, nil)

	return string(b), err
}
