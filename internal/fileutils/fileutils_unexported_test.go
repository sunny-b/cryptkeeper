package fileutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllParentDirectories(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name     string
		path     string
		expected []string
	}{
		{"Root directory", "/", []string{"/"}},
		{"Nested directory", "/a/b/c", []string{"/a/b/c", "/a/b", "/a", "/"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getAllParentDirectories(tt.path)
			assert.Equal(tt.expected, result)
		})
	}
}
func TestFindUp(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name      string
		searchDir string
		filename  string
		expected  string
	}{
		{"File found", "/", "testfile.txt", "/testfile.txt"},
		{"Nested file found", "/dir/foo/bar", "testfile2.txt", "/dir/testfile2.txt"},
		{"File not found", "/dir/foo/bar", "nonexistent.txt", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findUp(tt.searchDir, tt.filename)
			assert.Equal(tt.expected, result)
		})
	}
}
