package shell_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sunny-b/cryptkeeper/internal/shell"
)

func TestBashEscape(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", "''"},
		{"Literal characters", "abc", "$'abc'"},
		{"Single quote", "'", `$'\''`},
		{"Tab character", "\t", `$'\t'`},
		{"Newline character", "\n", `$'\n'`},
		{"Carriage return", "\r", `$'\r'`},
		{"Special characters", "&", "$'&'"},
		{"Hex character", string([]byte{6}), `$'\x06'`},
		{"Multiple special characters", "\t\n\r", `$'\t\n\r'`},
		{"Combination", "a\tb", "$'a\\tb'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shell.BashEscape(tt.input)
			assert.Equal(tt.expected, result)
		})
	}
}
