package shell

import (
	"path/filepath"
	"strings"
)

// Shell is the interface that represents the interaction with the host shell.
type Shell interface {
	// Hook is the string that gets evaluated into the host shell config and
	// setups direnv as a prompt hook.
	Hook() string

	// Export outputs the a string that exports the given environment variables
	Export(key, value string) string

	// ExportAll outputs a string that exports all the given environment
	ExportAll(Export) string

	// Unset unsets the given key from the host shell
	Unset(key string) string

	// Shell returns the name of the shell
	Shell() string

	// RCFile returns the path to the RC file
	RCFile() string
}

// DetectShell returns a Shell instance from the given target.
//
// target is usually $0 and can also be prefixed by `-`
func Detect(target string) Shell {
	target = filepath.Base(strings.ToLower(target))
	// $0 starts with "-"
	if target[0:1] == "-" {
		target = target[1:]
	}

	switch target {
	case "bash":
		return Bash
	case "fish":
		return Fish
	case "zsh":
		return Zsh
	}

	return nil
}
