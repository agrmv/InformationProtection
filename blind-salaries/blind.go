// Note that there are two types of signing/verifying going on.
// 1. The blind signature done by the Signer, done by SignSallary and
//    VerifySallary.  This uses rsablind.
// 2. The signing done by employees to ensure the authenticity of their message
//    to the Signer. Done by SignPSS/VerifyPSS.

package main

import (
	"../blind-salaries/fdh"
	"../blind-salaries/rsablind"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
)

var Keysize = 2048
var Hashize = 1536
var err error

/*********************
	   Utilities
**********************/
// verify blinded sig on the sallary
func VerifySallary(message, sig []byte, signerspubkey *rsa.PublicKey) error {
	hashed := fdh.Sum(crypto.SHA256, Hashize, message)
	return rsablind.VerifyBlindSignature(signerspubkey, hashed, sig)
}

// sign blinded message
func SignPSS(message []byte, privkey *rsa.PrivateKey) ([]byte, error) {
	rng := rand.Reader
	hashed := sha256.Sum256(message)
	return rsa.SignPSS(rng, privkey, crypto.SHA256, hashed[:], nil)
}

// verify sig on blinded message
func VerifyPSS(message, sig []byte, pubkey *rsa.PublicKey) error {
	hashed := sha256.Sum256(message)
	err = rsa.VerifyPSS(pubkey, crypto.SHA256, hashed[:], sig, nil)
	return err
}

func FinalMessage(salary, sig []byte) string {
	out := "salary: "
	out += string(salary) + "\nsignature: "
	out += base64.StdEncoding.EncodeToString(sig)
	return out
}

type BlindedMessage struct {
	Blinded   []byte
	Sig       []byte
	PublicKey rsa.PublicKey
}

/*********************
	   Voter
**********************/
type Voter struct {
	key         *rsa.PrivateKey
	signerskey  *rsa.PublicKey
	message     []byte
	unblinder   []byte
	PublicKey   *rsa.PublicKey
	has_blinded bool
}

func NewVoter(signerskey *rsa.PublicKey) (*Voter, error) {
	key, err := rsa.GenerateKey(rand.Reader, Keysize)
	if err != nil {
		return nil, errors.New("Error creating Voter RSA key")
	}
	return &Voter{key, signerskey, nil, nil, &key.PublicKey, false}, nil
}

// employee salary blinding function
func (e *Voter) BlindSalary(message []byte) (*BlindedMessage, error) {
	if e.has_blinded {
		return nil, errors.New("Voter already blinded a message")
	} else {
		e.has_blinded = true
	}

	e.message = message

	// We do a SHA256 full-domain-hash expanded to 1536 bits (3/4 the key size)
	hashed := fdh.Sum(crypto.SHA256, Hashize, message)

	// Blind the hashed message
	blind, unblinder, err := rsablind.Blind(e.signerskey, hashed)
	if err != nil {
		return nil, errors.New("Error blinding message")
	}
	e.unblinder = unblinder

	sig, err := SignPSS(blind, e.key)
	if err != nil {
		return nil, errors.New("Error signing blinded message")
	}
	return &BlindedMessage{blind, sig, *e.PublicKey}, nil
}

// employee unblinds and checks sig
func (e *Voter) Unblind(blindSig []byte) ([]byte, error) {
	// Unblind the signature
	unBlindedSig := rsablind.Unblind(e.signerskey, blindSig, e.unblinder)

	// verify the sig
	err = e.VerifySallary(e.message, unBlindedSig, e.signerskey)
	if err != nil {
		return nil, err
	}
	return unBlindedSig, nil
}

func (e *Voter) VerifySallary(message, sig []byte, signerspubkey *rsa.PublicKey) error {
	return VerifySallary(message, sig, e.signerskey)
}

/*********************
	   Signer
**********************/
type Signer struct {
	key       *rsa.PrivateKey
	employees map[rsa.PublicKey]bool
	PublicKey *rsa.PublicKey
}

func NewSigner() (*Signer, error) {
	key, err := rsa.GenerateKey(rand.Reader, Keysize)
	if err != nil {
		return nil, err
	}
	return &Signer{key, nil, &key.PublicKey}, nil
}

func (s *Signer) AddVoter(pubkeys []rsa.PublicKey) {
	s.employees = make(map[rsa.PublicKey]bool, len(pubkeys))

	for _, pk := range pubkeys {
		s.employees[pk] = false
	}
}

func (s *Signer) SignSalary(message *BlindedMessage) (sig []byte, err error) {
	err = s.authSignature(message)
	if err != nil {
		return nil, err
	}
	sig, err = rsablind.BlindSign(s.key, message.Blinded)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

func (s *Signer) authSignature(message *BlindedMessage) error {
	val, ok := s.employees[message.PublicKey]
	if !ok {
		return errors.New("Voter not registered")
	}
	if val {
		return errors.New("Voter already sent message")
	}
	return VerifyPSS(message.Blinded, message.Sig, &message.PublicKey)
}

func main() {
	// Установить подписавшего.
	signer, _ := NewSigner()

	// Инициализция избирателей.
	alice, _ := NewVoter(signer.PublicKey)
	bob, _ := NewVoter(signer.PublicKey)

	// Подписать избирателей.
	signer.AddVoter([]rsa.PublicKey{*alice.PublicKey, *bob.PublicKey})

	// Алиса и Боб записывают свой голос
	aliceSal := []byte("ZA")
	bobSal := []byte("PROTIV")

	// Они "ослепляют" свои голоса и подписывают их
	aliceBlindMsg, _ := alice.BlindSalary(aliceSal)
	bobBlindMsg, _ := bob.BlindSalary(bobSal)

	// Подписывающая сторона проверяет, что каждое сообщение авторизовано, и подписывает его.
	// Сообщение зашифровано (ослеплено). Подписывающий не может прочитать это.
	aliceBlindSig, _ := signer.SignSalary(aliceBlindMsg)
	bobBlindSig, _ := signer.SignSalary(bobBlindMsg)

	// Алиса и Боб снимают с себя подпись.
	aliceSig, _ := alice.Unblind(aliceBlindSig)
	bobSig, _ := bob.Unblind(bobBlindSig)

	// Они размещают свои голоса и подписи где-то анонимно.
	fmt.Println(FinalMessage(aliceSal, aliceSig))
	fmt.Println(FinalMessage(bobSal, bobSig))
}
