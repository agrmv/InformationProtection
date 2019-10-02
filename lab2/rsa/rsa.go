//package rsa
package main

import (
	"../../methods"
	"../renameMe"
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
	file, fileSize := renameMe.ReadFile("lab2/resourcesGlobal/test.jpg")

	encryptMessage := make([]int64, fileSize)
	decryptMessage := make([]byte, fileSize)
	encryptMessageBytes := make([]byte, fileSize)

	for i, v := range file {
		encryptMessage[i] = Encrypt(keys.publicKey, int64(v))
	}

	for i, v := range encryptMessage {
		decryptMessage[i] = byte(Decrypt(keys.privateKey, v))
	}

	for i, v := range encryptMessage {
		encryptMessageBytes[i] = byte(v)
	}

	renameMe.WriteFile("lab2/rsa/resources/encode.jpg", encryptMessageBytes)
	renameMe.WriteFile("lab2/rsa/resources/decode.jpg", decryptMessage)
}
