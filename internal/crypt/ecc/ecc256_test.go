package ecc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sunny-b/cryptkeeper/internal/crypt/ecc"
)

func TestEncryptDecrypt(t *testing.T) {
	assert := assert.New(t)
	e := &ecc.ECC256{}

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
			k, err := ecc.EphermalKey()
			assert.NoError(err)

			// Test encryption
			ciphertext, err := e.Encrypt(tt.plaintext, k)
			assert.NoError(err)
			assert.NotEqual(tt.plaintext, ciphertext)

			// Test decryption
			decrypted, err := e.Decrypt(ciphertext, k)
			assert.NoError(err)
			assert.Equal(tt.plaintext, string(decrypted))
		})
	}
}
