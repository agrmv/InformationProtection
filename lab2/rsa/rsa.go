//package rsa
package main

import (
	"../../lab1/euclid"
	"../../lab1/modularPow"
	"fmt"
	"math/rand"
)

type Pair struct {
	first, second int64
}

type Keys struct {
	publicKey  Pair
	privateKey Pair
}

//https://ru.wikipedia.org/wiki/%D0%A2%D0%B5%D1%81%D1%82_%D0%9C%D0%B8%D0%BB%D0%BB%D0%B5%D1%80%D0%B0_%E2%80%94_%D0%A0%D0%B0%D0%B1%D0%B8%D0%BD%D0%B0
func rabinMiller(n int64) uint8 {
	var ok uint8 = 1
	for i := 1; i <= 5 && ok == 1; i++ {
		var a = rand.Int63n(n-2) + 2
		var result = modularPow.ModularPow(a, n-1, n)
		if result == 1 || result == n-1 {
			ok &= 1
		}
	}
	return ok
}

func generatePrimaryKey() int64 {
	generated := rand.Int63() % 1000
	for rabinMiller(generated) == 0 {
		generated = rand.Int63()
	}
	return generated
}

func generatePublicKey(n int64) int64 {
	generated := rand.Int63() % 1000
	for gcd, _, _ := euclid.Gcd(n, generated); gcd != 1; {
		generated = rand.Int63() % 1000
	}
	return generated
}

func GCDRecursive(p, q int64) int64 {
	if q == 0 {
		return p
	}

	r := p % q
	return GCDRecursive(q, r)
}

func modularInverse(n int64, mod int64) int64 {
	_, _, inverse := euclid.Gcd(n, mod)
	fmt.Println(inverse)

	for inverse < 0 {
		inverse += mod
	}
	return inverse
}

//e открытый ключ
//d закрытый ключ
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

	result.publicKey = Pair{n, e}

	//вычисление секретной экспаненты
	d := modularInverse(e, phi)
	result.privateKey = Pair{n, d}
	return result
}

func Encrypt(key Keys, value int64) int64 {
	return modularPow.ModularPow(value, key.publicKey.second, key.publicKey.first)
}

func Decrypt(key Keys, value int64) int64 {
	return modularPow.ModularPow(value, key.privateKey.second, key.privateKey.first)
}

func main() {
	keys := generateKeys()
	fmt.Printf("Public key: %d, %d\n", keys.publicKey.first, keys.publicKey.second)
	fmt.Printf("Private key: %d, %d\n", keys.privateKey.first, keys.privateKey.second)

}
