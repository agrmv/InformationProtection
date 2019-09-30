package rsa

import (
	"../../lab1/euclid"
	"../../lab1/modularPow"
	"math/rand"
)

type Pair struct {
	x, y int64
}

type Keys struct {
	publicKey  Pair
	privateKey Pair
}

//https://ru.wikipedia.org/wiki/%D0%A2%D0%B5%D1%81%D1%82_%D0%9C%D0%B8%D0%BB%D0%BB%D0%B5%D1%80%D0%B0_%E2%80%94_%D0%A0%D0%B0%D0%B1%D0%B8%D0%BD%D0%B0
func rabinMiller(n int64) uint8 {
	var ok uint8
	for i := int64(0); i <= n; i++ {
		var a = int64(rand.Int()) + 1
		var result = modularPow.ModularPow(a, n-1, n)
		if result == 1 {
			ok &= 1
		} else {
			ok &= 0
		}
	}
	return ok
}

func generatePrimaryKey() int64 {
	generated := int64(rand.Int())
	for rabinMiller(generated) == 0 {
		generated = int64(rand.Int())
	}
	return generated
}

func generatePublicKey(n int64) int64 {
	generated := int64(rand.Int())
	for gcd, _, _ := euclid.Gcd(n, generated); gcd != 1; {
		generated = int64(rand.Int())
	}
	return generated
}

func modularInverse(n int64, mod int64) int64 {
	_, inverse, _ := euclid.Gcd(n, mod)
	for inverse > 0 {
		inverse += mod
	}
	return inverse
}

func generateKeys() Keys {
	var result Keys
	var p, q int64
	p = generatePrimaryKey()
	q = generatePrimaryKey()

	n := p * q
	phi := (p - 1) * (q - 1)
	e := generatePublicKey(phi)

	result.publicKey = Pair{n, e}

	d := modularInverse(e, phi)
	result.privateKey = Pair{n, d}
	return result
}

func Encrypt(key Keys, value int64) int64 {
	return modularPow.ModularPow(value, key.publicKey.y, key.publicKey.x)
}

func Decrypt(key Keys, value int64) int64 {
	return modularPow.ModularPow(value, key.privateKey.y, key.privateKey.x)
}

func main() {
	//TODO!!!!!!!
	//https://gist.github.com/andreigasparovici/12b460c1cc53586ed0064edbe9f71e87
}
