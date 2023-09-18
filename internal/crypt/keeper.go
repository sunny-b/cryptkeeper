package crypt

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/afero"
	"github.com/sunny-b/cryptkeeper/internal/crypt/aes"
	"github.com/sunny-b/cryptkeeper/internal/crypt/ecc"
	"github.com/sunny-b/cryptkeeper/internal/crypt/rsa"
	"github.com/sunny-b/cryptkeeper/internal/crypt/serpent"
	"github.com/sunny-b/cryptkeeper/internal/fileutils"
)

var fs afero.Fs = afero.NewOsFs()

type Encrypter interface {
	Encrypt(plaintext string, key any) (string, error)
	Decrypt(cipherText string, key any) (string, error)
}

type Keeper struct {
	encryptionType EncryptionType
	keyPath        string

	encrypter     Encrypter
	encryptionKey any
}

func NewKeeper(t EncryptionType, keyPath string) (*Keeper, error) {
	k := &Keeper{
		encryptionType: t,
		keyPath:        keyPath,
	}

	err := k.lazyInit()
	if err != nil {
		return nil, err
	}

	return k, nil
}

func GenerateKeys(enc EncryptionType, keyPath string) error {
	if err := validateEncryptionType(enc); err != nil {
		return err
	}

	var keys any
	var err error
	switch enc {
	case AES256:
		keys, err = aes.GenerateKeys()
	case ECC256:
		keys, err = ecc.GenerateKeys()
	case RSA2048:
		keys, err = rsa.GenerateKeys()
	case Serpent256:
		keys, err = serpent.GenerateKeys()
	default:
		err = ErrUnknownEncryptionType
	}
	if err != nil {
		return err
	}

	err = saveKeys(keys, keyPath)
	if err != nil {
		return fmt.Errorf("failed to write encryption key: %w", err)
	}

	return err
}

func (k *Keeper) Encrypt(secretName, plainText string) (string, error) {
	if err := k.lazyInit(); err != nil {
		return "", err
	}

	switch k.encryptionType {
	case AES256, RSA2048, Serpent256:
		return k.encrypter.Encrypt(plainText, k.encryptionKey)
	}

	keys, ok := k.encryptionKey.(*ecc.Keys)
	if !ok {
		return "", errors.New("corrupted key file")
	}

	key, err := ecc.EphermalKey()
	if err != nil {
		return "", errors.New("failed to ephermal key")
	}

	cipher, err := k.encrypter.Encrypt(plainText, key)
	if err != nil {
		return "", err
	}

	keys.KeyMap[secretName] = key

	err = saveKeys(keys, k.keyPath)
	if err != nil {
		return "", err
	}

	return cipher, nil
}

func (k *Keeper) Decrypt(secretName, cipher string) (string, error) {
	if err := k.lazyInit(); err != nil {
		return "", err
	}

	switch k.encryptionType {
	case AES256, RSA2048, Serpent256:
		return k.encrypter.Decrypt(cipher, k.encryptionKey)
	}

	keys, ok := k.encryptionKey.(*ecc.Keys)
	if !ok {
		return "", errors.New("corrupted key file")
	}

	key, ok := keys.KeyMap[secretName]
	if !ok {
		return "", errors.New("failed to find key for secret")
	}

	return k.encrypter.Decrypt(cipher, key)
}

func (k *Keeper) RemoveKey(secretName string) error {
	if err := k.lazyInit(); err != nil {
		return err
	}

	switch k.encryptionType {
	case AES256, RSA2048, Serpent256:
		return nil
	}

	keys, ok := k.encryptionKey.(*ecc.Keys)
	if !ok {
		return errors.New("corrupted key file")
	}

	delete(keys.KeyMap, secretName)

	return saveKeys(keys, k.keyPath)
}

func (k *Keeper) lazyInit() error {
	if err := validateEncryptionType(k.encryptionType); err != nil {
		return err
	}
	if k.encrypter == nil {
		err := k.fetchEncrypter()
		if err != nil {
			return err
		}
	}
	if k.encryptionKey == nil {
		err := k.fetchKeys()
		if err != nil {
			return err
		}
	}

	return nil
}
func saveKeys(keys any, keyPath string) error {
	b, err := json.Marshal(keys)
	if err != nil {
		return err
	}

	return fileutils.WriteFile(keyPath, b, 0644)
}

func (k *Keeper) fetchEncrypter() error {
	switch k.encryptionType {
	case AES256:
		k.encrypter = new(aes.AES256)
	case ECC256:
		k.encrypter = new(ecc.ECC256)
	case RSA2048:
		k.encrypter = new(rsa.RSA2048)
	case Serpent256:
		k.encrypter = new(serpent.Serpent256)
	}

	return nil
}

func (k *Keeper) fetchKeys() error {
	b, err := afero.ReadFile(fs, k.keyPath)
	if err != nil {
		return err
	}

	var key any
	switch k.encryptionType {
	case AES256:
		key = new(aes.EncryptionKey)
	case ECC256:
		key = new(ecc.Keys)
	case RSA2048:
		key = new(rsa.Keys)
	case Serpent256:
		key = new(serpent.EncryptionKey)
	}

	err = json.Unmarshal(b, key)
	if err != nil {
		return err
	}

	k.encryptionKey = key

	return nil
}

// func (k *Keeper) saveKeys(keys any) error {
// 	b, err := json.Marshal(keys)
// 	if err != nil {
// 		return err
// 	}

// 	return afero.WriteFile(fs, k.keyPath, b, 0644)
// }

func validateEncryptionType(t EncryptionType) error {
	switch t {
	case AES256, ECC256, RSA2048, Serpent256:
		return nil
	default:
		return ErrUnknownEncryptionType
	}
}
