package ecc_test

import (
	"crypto/elliptic"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sunny-b/cryptkeeper/internal/crypt/ecc"
)

func TestCustomJSONMarshalingAndUnmarshaling(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name  string
		curve elliptic.Curve
	}{
		{"P-256", elliptic.P256()},
		{"P-384", elliptic.P384()},
		{"P-521", elliptic.P521()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate ECC keys
			keys, err := ecc.EphermalKey()
			assert.NoError(err)

			// Marshal to JSON
			jsonData, err := json.Marshal(keys)
			assert.NoError(err)

			tmp := make(map[string]string)

			err = json.Unmarshal(jsonData, &tmp)
			assert.NoError(err)

			assert.NotEmpty(tmp["private_key"])
			assert.NotEmpty(tmp["public_key"])

			// Unmarshal from JSON
			newKeys := &ecc.Key{}
			err = json.Unmarshal(jsonData, newKeys)
			assert.NoError(err)

			// Validate that the keys are the same
			assert.Equal(keys.Private.D, newKeys.Private.D)
			assert.Equal(keys.Private.X, newKeys.Private.X)
			assert.Equal(keys.Private.Y, newKeys.Private.Y)
		})
	}
}
