package main

import (
	"../../ciphers"
	"../../methods"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
)

type rsaPublicKey struct {
	*rsa.PublicKey
}

type rsaPrivateKey struct {
	*rsa.PrivateKey
}

// loadPrivateKey loads an parses a PEM encoded private key file.
func loadPublicKey(publicKey *rsa.PublicKey) (Unsigner, error) {
	return parsePublicKey(ciphers.PublicKeyToBytes(publicKey))
}

// parsePublicKey parses a PEM encoded private key.
func parsePublicKey(pemBytes []byte) (Unsigner, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("ssh: no key found")
	}

	var rawkey interface{}
	switch block.Type {
	case "PUBLIC KEY":
		rsa, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rawkey = rsa
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	}

	return newUnsignerFromKey(rawkey)
}

// loadPrivateKey loads an parses a PEM encoded private key file.
func loadPrivateKey(privateKey *rsa.PrivateKey) (Signer, error) {
	return parsePrivateKey(ciphers.PrivateKeyToBytes(privateKey))
}

// parsePublicKey parses a PEM encoded private key.
func parsePrivateKey(pemBytes []byte) (Signer, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("ssh: no key found")
	}

	var rawkey interface{}
	switch block.Type {
	case "RSA PRIVATE KEY":
		rsa, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rawkey = rsa
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	}
	return newSignerFromKey(rawkey)
}

// A Signer is can create signatures that verify against a public key.
type Signer interface {
	// Sign returns raw signature for the given data. This method
	// will apply the hash specified for the keytype to the data.
	Sign(data []byte) ([]byte, error)
}

// A Signer is can create signatures that verify against a public key.
type Unsigner interface {
	// Sign returns raw signature for the given data. This method
	// will apply the hash specified for the keytype to the data.
	Unsign(data []byte, sig []byte) error
}

func newSignerFromKey(k interface{}) (Signer, error) {
	var sshKey Signer
	switch t := k.(type) {
	case *rsa.PrivateKey:
		sshKey = &rsaPrivateKey{t}
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %T", k)
	}
	return sshKey, nil
}

func newUnsignerFromKey(k interface{}) (Unsigner, error) {
	var sshKey Unsigner
	switch t := k.(type) {
	case *rsa.PublicKey:
		sshKey = &rsaPublicKey{t}
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %T", k)
	}
	return sshKey, nil
}

// Sign signs data with rsa-sha256
func (r *rsaPrivateKey) Sign(data []byte) ([]byte, error) {
	h := sha256.New()
	h.Write(data)
	d := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, crypto.SHA256, d)
}

// Unsign verifies the message using a rsa-sha256 signature
func (r *rsaPublicKey) Unsign(message []byte, sig []byte) error {
	h := sha256.New()
	h.Write(message)
	d := h.Sum(nil)
	return rsa.VerifyPKCS1v15(r.PublicKey, crypto.SHA256, d, sig)
}

func main() {
	fileToSign, fileSize := methods.ReadFile("lab2/resourcesGlobal/test.jpg")
	privateKey, publicKey := ciphers.GenerateKeyPair(fileSize)

	signer, err := loadPrivateKey(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	signed, err := signer.Sign(fileToSign)
	if err != nil {
		log.Fatal(err)
	}
	//можно спрятать в функцию
	sig := base64.StdEncoding.EncodeToString(signed)
	//тут записывать в файл
	fmt.Printf("Signature: %v\n", sig)

	parser, perr := loadPublicKey(publicKey)
	if perr != nil {
		log.Fatal(err)
	}

	err = parser.Unsign(fileToSign, signed)
	if err != nil {
		log.Fatal(err)
	}

	/*TODO:
	 * Записывать в файл его подпись(на бд сделаю)
	 * Переписать rsa на crypto DONE
	 * Переписать генерацию ключей под pem DONE
	 * Навести порядок с файловой структурой, по сути сделать аналог crypro, с некоторыми самописными функциями
	 * Возможно стоит юзать не sha256
	 * Вынести структуры с интерфейсами в родительский пакет???
	 * elgamal по аналогии должно быть легко
	 * c гостом разобраться че это такое и просто наебать
	 */
	//fmt.Printf("Unsign error: %v\n", err)
}
