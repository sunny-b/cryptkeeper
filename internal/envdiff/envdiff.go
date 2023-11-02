package envdiff

import (
	"os"
	"strings"

	"github.com/direnv/direnv/v2/gzenv"
	"github.com/sunny-b/cryptkeeper/internal/config"
	"github.com/sunny-b/cryptkeeper/internal/crypt"
	"github.com/sunny-b/cryptkeeper/internal/shell"
	"github.com/sunny-b/cryptkeeper/internal/utils"
)

// Copied from Direnv to handle env diffing

// IgnoredKeys is list of keys we don't want to deal with
var IgnoredKeys = map[string]bool{
	"COMP_WORDBREAKS": true, // Avoids segfaults in bash
	"PS1":             true, // PS1 should not be exported, fixes problem in bash

	// variables that should change freely
	"OLDPWD":    true,
	"PWD":       true,
	"SHELL":     true,
	"SHELLOPTS": true,
	"SHLVL":     true,
	"_":         true,
}

// EnvDiff represents the diff between two environments
type EnvDiff struct {
	Revert map[string]string
	Last   map[string]string
}

// NewEnvDiff is an empty constructor for EnvDiff
func NewEnvDiff() *EnvDiff {
	return &EnvDiff{make(map[string]string), make(map[string]string)}
}

// BuildEnvDiff analyses the changes between 'e1' and 'e2' and builds an
// EnvDiff out of it.
func BuildEnvDiff(e1, e2 config.Env) *EnvDiff {
	if e1 == nil && e2 == nil {
		return nil
	}

	diff := NewEnvDiff()

	for key := range e1 {
		if IgnoredEnv(key) {
			continue
		}
		if e2[key] != e1[key] || !utils.In(key, e2) {
			diff.Revert[key] = e1[key]
		}
	}

	for key := range e2 {
		if IgnoredEnv(key) {
			continue
		}
		if e2[key] != e1[key] || !utils.In(key, e1) {
			diff.Last[key] = e2[key]
		}
	}

	return diff
}

// LoadEnvDiff unmarshalls a gzenv string back into an EnvDiff.
func LoadEnvDiff(gzenvStr string) (diff *EnvDiff, err error) {
	diff = new(EnvDiff)
	err = gzenv.Unmarshal(gzenvStr, diff)
	return
}

// Any returns if the diff contains any changes.
func (d *EnvDiff) Any() bool {
	return len(d.Revert) > 0 || len(d.Last) > 0
}

// ToShell applies the env diff as a set of commands that are understood by
// the target `shell`. The outputted string is then meant to be evaluated in
// the target shell.
func (d *EnvDiff) ToShell(sh shell.Shell) string {
	e := make(shell.Export)

	for key, value := range d.Revert {
		_, ok := d.Last[key]
		if !ok {
			if value == "" {
				e.Remove(key)
			} else {
				e.Add(key, value)
			}

		}
	}

	for key, value := range d.Last {
		e.Add(key, value)
	}

	return sh.ExportAll(e)
}

// Patch applies the diff to the given env and returns a new env with the
// changes applied.
func (d *EnvDiff) Patch(env config.Env) (newEnv config.Env) {
	newEnv = make(config.Env)

	for k, v := range env {
		newEnv[k] = v
	}

	for key := range d.Revert {
		delete(newEnv, key)
	}

	for key, value := range d.Last {
		newEnv[key] = value
	}

	return newEnv
}

// Reverse flips the diff so that it applies the other way around.
func (d *EnvDiff) Reverse() *EnvDiff {
	return &EnvDiff{Revert: d.Last, Last: d.Revert}
}

// Serialize marshalls the environment diff to the gzenv format.
func (d *EnvDiff) Serialize() string {
	return gzenv.Marshal(d)
}

func (d *EnvDiff) Decrypt(decrypter *crypt.Keeper) error {
	for key, value := range d.Revert {
		decrypted, err := decrypter.Decrypt(key, value)
		if err != nil {
			return err
		}

		d.Revert[key] = decrypted
	}

	for key, value := range d.Last {
		decrypted, err := decrypter.Decrypt(key, value)
		if err != nil {
			return err
		}

		d.Last[key] = decrypted
	}

	return nil
}

//// Utils

// IgnoredEnv returns true if the key should be ignored in environment diffs.
func IgnoredEnv(key string) bool {
	if strings.HasPrefix(key, "__fish") {
		return true
	}
	if strings.HasPrefix(key, "BASH_FUNC_") {
		return true
	}
	_, found := IgnoredKeys[key]
	return found
}

// FetchRevert undoes the recorded changes (if any) to the supplied environment,
// returning a new environment
func FetchRevert(envKey string) (map[string]*string, error) {
	env := make(map[string]*string)
	err := FetchEnv(envKey, &env)
	if err != nil {
		return env, err
	}

	return env, nil
}

func FetchLastEnv(diffKey string, keeper *crypt.Keeper) (config.Env, error) {
	env := make(config.Env)
	err := FetchEncryptedEnv(diffKey, keeper, &env)
	if err != nil {
		return env, err
	}

	return env, nil
}

func FetchEnv(envKey string, obj any) error {
	env, ok := os.LookupEnv(envKey)
	if !ok || len(env) == 0 {
		return nil
	}

	return LoadEnv(env, obj)
}

func FetchEncryptedEnv(envKey string, keeper *crypt.Keeper, obj any) error {
	env, ok := os.LookupEnv(envKey)
	if !ok || len(env) == 0 {
		return nil
	}

	decrypted, err := keeper.Decrypt(envKey, env)
	if err != nil {
		return err
	}

	return LoadEnv(decrypted, obj)
}

// LoadEnv unmarshalls a gzenv string back into an EnvDiff.
func LoadEnv(gzenvStr string, data any) (err error) {
	err = gzenv.Unmarshal(gzenvStr, data)
	return
}

// func FetchDiff(revertKey, lastKey string, keeper *crypt.Keeper) (*EnvDiff, error) {
// 	revert, err := FetchRevert(revertKey)
// 	if err != nil {
// 		return nil, err
// 	}

// 	last, err := FetchLastEnv(lastKey, keeper)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return BuildEnvDiff(revert, last), nil
// }

// func osEnv() config.Env {
// 	env := make(config.Env)

// 	for _, kv := range os.Environ() {
// 		kv2 := strings.SplitN(kv, "=", 2)

// 		key := kv2[0]
// 		value := kv2[1]

// 		env[key] = value
// 	}

// 	return env
// }
