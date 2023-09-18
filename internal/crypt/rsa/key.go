package rsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
)

//nolint:musttag
type Keys struct {
	Private *rsa.PrivateKey
}

func (k *Keys) MarshalJSON() ([]byte, error) {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(k.Private)
	privateKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		},
	)

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&k.Private.PublicKey)
	if err != nil {
		return nil, err
	}
	publicKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		},
	)

	return json.Marshal(map[string]string{
		"private_key": base64.StdEncoding.EncodeToString(privateKeyPem),
		"public_key":  base64.StdEncoding.EncodeToString(publicKeyPem),
	})
}

func (k *Keys) UnmarshalJSON(data []byte) error {
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
	k.Private, err = x509.ParsePKCS1PrivateKey(block.Bytes)
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
	k.Private.PublicKey = *(pub.(*rsa.PublicKey))

	return nil
}
