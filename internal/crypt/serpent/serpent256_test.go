package serpent_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sunny-b/cryptkeeper/internal/crypt/serpent"
)

func TestEncryptDecrypt(t *testing.T) {
	assert := assert.New(t)
	r := &serpent.Serpent256{}

	keys, err := serpent.GenerateKeys()
	assert.NoError(err)

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
			ciphertext, err := r.Encrypt(tt.plaintext, keys)
			assert.NoError(err)
			assert.NotEqual(tt.plaintext, ciphertext)

			decrypted, err := r.Decrypt(ciphertext, keys)
			assert.NoError(err)
			assert.Equal(tt.plaintext, string(decrypted))
		})
	}
}
