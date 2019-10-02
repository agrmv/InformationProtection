//package rsa
package main

import (
	"../../methods"
	"fmt"
	"math/rand"
	"time"
)

type Pair struct {
	first, second int64
}

type Keys struct {
	publicKey  Pair
	privateKey Pair
}

func generatePrivateKey() int64 {
	generated := rand.Int63() % 1000
	for !methods.TestFerma(generated, 5) {
		generated = rand.Int63() % 1000
	}
	return generated
}

//(1 < e < phi)
func getOpenExp(phi int64) int64 {
	generated := rand.Int63n(phi-1) + 1
	for !methods.TestFerma(generated, 5) {
		generated = rand.Int63n(phi-1) + 1
	}
	return generated
}

//d = e^-1 mod phi
func getSecretExp(e, phi int64) int64 {
	_, inverse, _ := methods.GcdExtended(e, phi)
	for inverse < 0 {
		inverse += phi
	}
	return inverse
}

//https://ru.wikipedia.org/wiki/RSA
func generateKeys() Keys {
	var result Keys
	var p, q int64

	p = generatePrivateKey()
	q = generatePrivateKey()

	n := p * q
	// вычисление функции эйлера
	phi := (p - 1) * (q - 1)
	e := getOpenExp(phi)

	result.publicKey = Pair{e, n}

	//вычисление секретной экспаненты
	d := getSecretExp(e, phi)

	result.privateKey = Pair{d, n}
	return result
}

func Encrypt(key Pair, value int64) int64 {
	return methods.ModularPow(value, key.first, key.second)
}

func Decrypt(key Pair, value int64) int64 {
	return methods.ModularPow(value, key.first, key.second)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	keys := generateKeys()
	fmt.Printf("Public key: %d, %d\n", keys.publicKey.first, keys.publicKey.second)
	fmt.Printf("Private key: %d, %d\n", keys.privateKey.first, keys.privateKey.second)
	message := "SLAVA UKRAINE"
	fmt.Println("Initial Message: " + message)

	encryptMessage := make([]int64, len(message))
	decryptMessage := make([]int64, len(message))

	for i := range message {
		encryptMessage[i] = Encrypt(keys.publicKey, int64(message[i]))
	}

	for i := range message {
		decryptMessage[i] = Decrypt(keys.privateKey, encryptMessage[i])
	}

	for i := range message {
		fmt.Print(string(decryptMessage[i]))
	}
}
