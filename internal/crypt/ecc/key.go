package ecc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
)

type Keys struct {
	KeyMap map[string]*Key `json:"keys"`
}

//nolint:musttag
type Key struct {
	Private *ecdsa.PrivateKey
}

func GenerateKeys() (Keys, error) {
	keys := Keys{
		KeyMap: make(map[string]*Key),
	}

	return keys, nil
}

func EphermalKey() (*Key, error) {
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return &Key{k}, nil
}

func (k *Key) MarshalJSON() ([]byte, error) {
	privateKeyBytes, err := x509.MarshalECPrivateKey(k.Private)
	if err != nil {
		return nil, err
	}

	privateKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&k.Private.PublicKey)
	if err != nil {
		return nil, err
	}

	publicKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return json.Marshal(map[string]string{
		"private_key": base64.StdEncoding.EncodeToString(privateKeyPem),
		"public_key":  base64.StdEncoding.EncodeToString(publicKeyPem),
	})
}

func (k *Key) UnmarshalJSON(data []byte) error {
	var m map[string]string
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	rawPrivate, err := base64.StdEncoding.DecodeString(m["private_key"])
	if err != nil {
		return err
	}

	block, _ := pem.Decode(rawPrivate)
	k.Private, err = x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	rawPublic, err := base64.StdEncoding.DecodeString(m["public_key"])
	if err != nil {
		return err
	}

	block, _ = pem.Decode(rawPublic)
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	k.Private.PublicKey = *(pub.(*ecdsa.PublicKey))

	return nil
}
