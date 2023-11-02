package fileutils_test

import (
	"errors"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/sunny-b/cryptkeeper/internal/fileutils" // Replace with the actual import path
)

var fs = afero.NewMemMapFs()

func TestMain(m *testing.M) {
	fileutils.SetFS(fs)

	// Create some test files and directories in the in-memory file system
	err := afero.WriteFile(fs, "/testfile.txt", []byte("Hello, world!"), 0644)
	if err != nil {
		panic(err)
	}

	err = fs.Mkdir("/dir", 0755)
	if err != nil {
		panic(err)
	}

	err = afero.WriteFile(fs, "/dir/testfile2.txt", []byte("Another file"), 0644)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestTextExistsInFile(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name       string
		filePath   string
		targetText string
		expected   bool
		err        error
	}{
		{"Text exists", "/testfile.txt", "Hello", true, nil},
		{"Text does not exist", "/testfile.txt", "Bye", false, nil},
		{"File does not exist", "/nonexistent.txt", "Hello", false, afero.ErrFileNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fileutils.TextExistsInFile(tt.filePath, tt.targetText)
			assert.Equal(tt.expected, result)
			if tt.err != nil && assert.Error(err) {
				assert.ErrorContains(err, tt.err.Error())
			}
		})
	}
}

func TestClean(t *testing.T) {
	assert := assert.New(t)

	wd, err := os.Getwd()
	assert.NoError(err)

	tests := []struct {
		name     string
		filePath string
		expected string
	}{
		{"Clean path", "/dir/../testfile.txt", "/testfile.txt"},
		{"Already clean", "/testfile.txt", "/testfile.txt"},
		{"Dot in front clean", "./testfile.txt", wd + "/testfile.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fileutils.Clean(tt.filePath)
			assert.Equal(tt.expected, result)
		})
	}
}

func TestFindPathTo(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name     string
		file     string
		expected string
		err      error
	}{
		{"File exists", "testfile.txt", "/testfile.txt", nil},
		{"File does not exist", "nonexistent.txt", "", errors.New("file not found")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fileutils.FindPathTo(tt.file)
			assert.Equal(tt.expected, result)
			assert.Equal(tt.err, err)
		})
	}
}

func TestFileExists(t *testing.T) {
	assert := assert.New(t)

	filePath := "/exists.txt"

	// Create a temporary file
	tempFile, err := fs.Create(filePath)
	assert.NoError(err)
	defer os.Remove(tempFile.Name())

	tests := []struct {
		name     string
		filePath string
		expected bool
	}{
		{"File exists", filePath, true},
		{"File does not exist", "nonexistent.txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fileutils.FileExists(tt.filePath)
			assert.Equal(tt.expected, result)
		})
	}
}

func TestIsChildDirOrSame(t *testing.T) {
	tests := []struct {
		name   string
		parent string
		child  string
		want   bool
	}{
		{
			name:   "IsChild",
			parent: "/home/user",
			child:  "/home/user/file.txt",
			want:   true,
		},
		{
			name:   "Same",
			parent: "/home/user",
			child:  "/home/user",
			want:   true,
		},
		{
			name:   "NotChild",
			parent: "/home/user",
			child:  "/home/another/file.txt",
			want:   false,
		},
		{
			name:   "EmptyPath",
			parent: "/home/user",
			child:  "",
			want:   false,
		},
		{
			name:   "ChildIsParent",
			parent: "/home/user",
			child:  "/home",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := fileutils.IsChildDirOrSame(tt.child, tt.parent)
			assert.Equal(t, tt.want, got)
		})
	}
}
