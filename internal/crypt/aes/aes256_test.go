package aes_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sunny-b/cryptkeeper/internal/crypt/aes"
)

func TestAES256(t *testing.T) {
	assert := assert.New(t)
	key, err := aes.GenerateKeys()
	assert.NoError(err)

	a := &aes.AES256{}

	// Table-driven tests
	tests := []struct {
		name      string
		plaintext string
	}{
		{"Normal text", "Hello, world!"},
		{"Empty text", ""},
		{"Special characters", "!@#$\n%^\t&*()"},
		{"Long text", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus lacinia odio vitae vestibulum."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encryption
			cipherText, err := a.Encrypt(tt.plaintext, key)
			assert.NoError(err)

			// Decryption
			decrypted, err := a.Decrypt(cipherText, key)
			assert.NoError(err)
			assert.Equal(tt.plaintext, string(decrypted))
		})
	}
}
