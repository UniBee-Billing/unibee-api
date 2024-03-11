package utility

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// RsaEncrypt encrypts data using rsa public key.
func RsaEncrypt(pubkey, data []byte) ([]byte, error) {
	block, _ := pem.Decode(pubkey)
	if block == nil {
		return nil, errors.New("decode public key error")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.EncryptPKCS1v15(rand.Reader, pub.(*rsa.PublicKey), data)
}

// RsaDecrypt decrypts data using rsa private key.
func RsaDecrypt(prvkey, cipher []byte) ([]byte, error) {
	block, _ := pem.Decode(prvkey)
	if block == nil {
		return nil, errors.New("decode private key error")
	}
	prv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, prv, cipher)
}
