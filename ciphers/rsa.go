package ciphers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"log"
)

type Signer interface {
	// Sign returns raw signature for the given data. This method
	// will apply the hash specified for the keytype to the data.
	Sign(data []byte) ([]byte, error)
}

type Unsigner interface {
	// Sign returns raw signature for the given data. This method
	// will apply the hash specified for the keytype to the data.
	Unsign(data []byte, sig []byte) error
}

// GenerateKeyPair generates a new key pair
func GenerateKeyPair(bits int64) (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, int(bits/100))
	if err != nil {
		log.Fatal(err)
	}
	return privateKey, &privateKey.PublicKey
}

// PrivateKeyToBytes private key to bytes
func PrivateKeyToBytes(privateKey *rsa.PrivateKey) []byte {
	privateBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)

	return privateBytes
}

// PublicKeyToBytes public key to bytes
func PublicKeyToBytes(publicKey *rsa.PublicKey) []byte {
	pubASN1, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		log.Fatal(err)
	}

	publicBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return publicBytes
}

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(privateKey []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(privateKey)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		log.Fatal(err)
	}
	return key
}

// BytesToPublicKey bytes to public key
func BytesToPublicKey(publicKey []byte) *rsa.PublicKey {
	block, _ := pem.Decode(publicKey)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		log.Fatal(err)
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		log.Fatal("not ok")
	}
	return key
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(message []byte, publicKey *rsa.PublicKey) []byte {
	hash := sha1.New()
	cipherMessage, err := rsa.EncryptOAEP(hash, rand.Reader, publicKey, message, nil)
	if err != nil {
		log.Fatal(err)
	}
	return cipherMessage
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(message []byte, privateKey *rsa.PrivateKey) []byte {
	hash := sha1.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, privateKey, message, nil)
	if err != nil {
		log.Fatal(err)
	}
	return plaintext
}
