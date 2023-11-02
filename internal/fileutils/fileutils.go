package fileutils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

var ErrFileNotFound = errors.New("file not found")

var fs = afero.NewOsFs()

func SetFS(f afero.Fs) {
	fs = f
}

func WriteFile(fileName string, data []byte, perm os.FileMode) error {
	return afero.WriteFile(fs, fileName, data, perm)
}

func ReadFile(fileName string) ([]byte, error) {
	return afero.ReadFile(fs, fileName)
}

func TextExistsInFile(filePath, targetText string) (bool, error) {
	content, err := afero.ReadFile(fs, filePath)
	if err != nil {
		return false, err
	}

	// Check if the target text exists in the file
	return strings.Contains(string(content), targetText), nil
}

func Clean(filePath string) string {
	filePath = filepath.Clean(filePath)

	absFile, err := filepath.Abs(filePath)
	if err != nil {
		return filePath
	}

	return absFile
}

func FindPathTo(file string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	path := findUp(wd, file)
	if path == "" {
		return "", ErrFileNotFound
	}

	return path, nil
}

// isChildDirOrSame determines if the path `child` is a child path or the same as `parent`.
func IsChildDirOrSame(child, parent string) (bool, error) {
	return strings.HasPrefix(Clean(child), Clean(parent)), nil
}

func findUp(searchDir, filename string) string {
	for _, dir := range getAllParentDirectories(searchDir) {
		path := filepath.Join(dir, filename)
		if FileExists(path) {
			return path
		}
	}

	return ""
}

func getAllParentDirectories(path string) (paths []string) {
	if path == "/" {
		return []string{"/"}
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return
	}

	paths = []string{}
	for absPath != "/" && absPath != "." {
		paths = append(paths, absPath)
		absPath = filepath.Dir(absPath)
	}

	paths = append(paths, "/")

	return
}

// copied from direnv code
func FileExists(path string) bool {
	// Some broken filesystems like SSHFS return file information on stat() but
	// then cannot open the file. So we use os.Open.
	f, err := fs.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	// Next, check that the file is a regular file.
	fi, err := f.Stat()
	if err != nil {
		return false
	}

	return fi.Mode().IsRegular()
}
