//package rsa
package main

import (
	"../../lab1/euclid"
	"../../lab1/modularPow"
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

func TestFerma(phi int64) bool {
	if phi == 2 {
		return true
	}
	if phi&1 != 1 {
		return false
	}
	for i := 0; i < 5; i++ {
		a := rand.Int63n(phi-1) + 1
		gcd, _, _ := euclid.Gcd(a, phi)
		if gcd != 1 || modularPow.ModularPow(a, phi-1, phi) != 1 {
			return false
		}
	}
	return true
}

func rabinMiller(n int64) bool {
	for i := 0; i < 5; i++ {
		a := rand.Int63n(n-2) + 2
		if modularPow.ModularPow(a, n-1, n) != 1 {
			return false
		}
	}
	return true
}

func generatePrimaryKey() int64 {
	generated := rand.Int63() % 1000
	for !rabinMiller(generated) {
		generated = rand.Int63() % 1000
	}
	return generated
}

func generatePublicKey(n int64) int64 {
	generated := rand.Int63n(n-1) + 1
	for !TestFerma(generated) {
		generated = rand.Int63n(n-1) + 1
	}
	return generated
}

func modularInverse(n int64, mod int64) int64 {
	_, inverse, _ := euclid.Gcd(n, mod)
	for ; inverse < 0; inverse += mod {
	}
	return inverse
}

//https://ru.wikipedia.org/wiki/RSA
func generateKeys() Keys {
	var result Keys
	var p, q int64

	p = generatePrimaryKey()
	q = generatePrimaryKey()

	n := p * q
	// вычисление функции эйлера
	phi := (p - 1) * (q - 1)
	e := generatePublicKey(phi)

	result.publicKey = Pair{e, n}

	//вычисление секретной экспаненты
	d := modularInverse(e, phi)

	result.privateKey = Pair{d, n}
	return result
}

func Encrypt(key Pair, value int64) int64 {
	return modularPow.ModularPow(value, key.first, key.second)
}

func Decrypt(key Pair, value int64) int64 {
	return modularPow.ModularPow(value, key.first, key.second)
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
