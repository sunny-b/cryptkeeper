package rsa_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sunny-b/cryptkeeper/internal/crypt/rsa"
)

func TestCustomJSONMarshalingAndUnmarshalingRSA(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name      string
		keyLength int
	}{
		{"2048 bits", 2048},
		{"3072 bits", 3072},
		{"4096 bits", 4096},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys, err := rsa.GenerateKeys()
			assert.NoError(err)

			jsonData, err := json.Marshal(keys)
			assert.NoError(err)

			tmp := make(map[string]string)

			err = json.Unmarshal(jsonData, &tmp)
			assert.NoError(err)

			assert.NotEmpty(tmp["private_key"])
			assert.NotEmpty(tmp["public_key"])

			newKeys := &rsa.Keys{}
			err = json.Unmarshal(jsonData, newKeys)
			assert.NoError(err)

			assert.Equal(keys.Private.D, newKeys.Private.D)
			assert.Equal(keys.Private.N, newKeys.Private.N)
		})
	}
}
