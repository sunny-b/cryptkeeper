package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/direnv/direnv/gzenv"

	"github.com/sunny-b/cryptkeeper/internal/crypt"
	"github.com/sunny-b/cryptkeeper/internal/fileutils"
)

type Mode string

const (
	envrc       = ".envrc"
	fileName    = ".ckrc"
	keyFileName = ".ckkey"

	CKWatchEnvKey  = "CK_WATCH"
	CKRevertEnvKey = "CK_REVERT"
	CKLastEnvKey   = "CK_LAST"

	DirenvMode     Mode = "direnv"
	StandaloneMode Mode = "standalone"
)

var (
	cachedConfig []byte
	cachedPath   string

	CKEnvKeys = []string{
		CKWatchEnvKey,
		CKRevertEnvKey,
		CKLastEnvKey,
	}
)

//nolint:musttag
type Config struct {
	Mode       Mode       `json:"mode"`
	Encryption Encryption `json:"encryption"`
	Env        Env        `json:"env"`

	Path string
}

func (c *Config) MarshalJSON() ([]byte, error) {
	tmp := struct {
		Mode       Mode       `json:"mode"`
		Encryption Encryption `json:"encryption"`
		Env        Env        `json:"env"`
	}{
		Mode:       c.Mode,
		Encryption: c.Encryption,
		Env:        c.Env,
	}

	return json.Marshal(tmp)
}

type Encryption struct {
	Type    crypt.EncryptionType `json:"type"`
	KeyPath string               `json:"key_path"`
}

type Direnv struct {
	RCPath string `json:"rc_path"`
}

func (d *Direnv) Enabled() bool {
	return d != nil && len(d.RCPath) > 0
}

type Env map[string]string

type AnyEnv interface {
	map[string]string | map[string]*string | Env
}

func (e Env) Keys() []string {
	keys := make([]string, 0, len(e))
	for key := range e {
		keys = append(keys, key)
	}
	return keys
}

func (e Env) Decrypt(decrypter *crypt.Keeper) error {
	for key, value := range e {
		decrypted, err := decrypter.Decrypt(key, value)
		if err != nil {
			return err
		}
		e[key] = decrypted
	}

	return nil
}

// Copy returns a fresh copy of the env. Because the env is a map under the
// hood, we want to get a copy whenever we mutate it and want to keep the
// original around.
func (e Env) Copy() Env {
	newEnv := make(Env)

	for key, value := range e {
		newEnv[key] = value
	}

	return newEnv
}

// Serialize marshalls the environment diff to the gzenv format.
func Serialize[E AnyEnv](e E) string {
	return gzenv.Marshal(e)
}

// CleanContext removes all the direnv-related environment variables. Call
// this after reverting the environment, otherwise direnv will just be amnesic
// about the previously-loaded environment.
func (e Env) CleanContext() {
	delete(e, CKLastEnvKey)
	delete(e, CKRevertEnvKey)
	delete(e, CKWatchEnvKey)
}

func (c *Config) IsDirenvIntegrated() bool {
	return c.Mode == DirenvMode
}

func FileName() string {
	return fileName
}

func KeyFileName() string {
	return keyFileName
}

func Path() (string, error) {
	return findPathToConfig()
}

func ReadInConfig() error {
	path, err := findPathToConfig()
	if err != nil {
		return err
	}

	b, err := fileutils.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	cachedConfig = b

	return nil
}

func Init(config *Config) error {
	path := config.Path

	// make sure the config file doesn't exist
	if fileutils.FileExists(path) {
		return nil
	}

	err := Write(config)
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	cachedPath = path

	return nil
}

func Write(config *Config) error {
	path := config.Path
	if path == "" {
		var err error
		path, err = findPathToConfig()
		if err != nil {
			return err
		}
	}

	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = fileutils.WriteFile(path, b, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	cachedConfig = b

	return nil
}

func GetConfigFromPath(path string) (*Config, error) {
	config := &Config{}
	err := loadConfig(path, config)
	if err != nil {
		return nil, err
	}

	config.Path = path

	return config, nil
}

func GetConfig() (*Config, error) {
	if len(cachedConfig) > 0 {
		return getCachedConfig()
	}

	path, err := findPathToConfig()
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = loadConfig(path, config)
	if err != nil {
		return nil, err
	}

	config.Path = path

	return config, nil
}

func findPathToConfig() (string, error) {
	if len(cachedPath) > 0 {
		return cachedPath, nil
	}

	nearestPath, err := fileutils.FindPathTo(fileName)
	if err != nil {
		// In the case where the user moves out of the directory, we want to
		// return the watch path so we can unload it.
		if watchPath, ok := os.LookupEnv(CKWatchEnvKey); ok {
			return watchPath, nil
		}

		return "", err
	}

	cachedPath = nearestPath

	return nearestPath, nil
}

func loadConfig(path string, config *Config) error {
	b, err := fileutils.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	err = json.Unmarshal(b, config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

func unmarshalConfig() (*Config, error) {
	if len(cachedConfig) == 0 {
		return nil, errors.New("no cached config")
	}

	config := &Config{}
	err := json.Unmarshal(cachedConfig, config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}

func getCachedConfig() (*Config, error) {
	config, err := unmarshalConfig()
	if err != nil {
		return nil, err
	}

	config.Path, err = findPathToConfig()
	if err != nil {
		return nil, err
	}

	return config, nil
}
