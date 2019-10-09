package main

import (
	"../../methods"
	"fmt"
	rand2 "math/rand"
)

type Pair struct {
	First  int64 `json:"First"`
	Second int64 `json:"Second"`
}

type Keys struct {
	AliceKeys Pair `json:"PublicKey"`
	BobKeys   Pair `json:"PrivateKey"`
}

/*TODO THIS
type Shamir struct {
	ShamirKeys Keys
	p          int64
	m          int64
}*/

func GenerateCD(keys *Pair, p int64) Pair {
	keys.First = int64(1)
	keys.Second = rand2.Int63n(p)
	for keys.First < p {
		if (keys.First*keys.Second)%(p-1) == 1 {
			return Pair{keys.First, keys.Second}
		}
		keys.First += 1
	}
	if keys.First == p {
		GenerateCD(keys, p)
	}
	return Pair{keys.First, keys.Second}
}

func generateAliceKeys(keys *Keys, p int64) {
	keys.AliceKeys = GenerateCD(&keys.AliceKeys, p) // генерируем взаимно-простые C и D для Алисы
}

func generateBobKeys(keys *Keys, p int64) {
	keys.BobKeys = GenerateCD(&keys.BobKeys, p) // генерируем взаимно-простые C и D для Боба
}

func Encrypt(alice Pair, m, p int64) int64 {
	return methods.ModularPow(methods.ModularPow(m, alice.First, p), alice.Second, p)
}

func Decrypt(bob Pair, p, encrypt int64) int64 {
	return methods.ModularPow(methods.ModularPow(encrypt, bob.First, p), bob.Second, p)
}

func main() {
	p := methods.DefaultGeneratePrime() // генерируем общее простое число
	keys := Keys{}
	generateAliceKeys(&keys, p)
	generateBobKeys(&keys, p)

	m := int64(45)

	encrypted := Encrypt(keys.AliceKeys, m, p)
	decrypted := Decrypt(keys.BobKeys, p, encrypted)

	fmt.Printf("Secret message: %d\n", m)

	fmt.Printf("step4: %d\n", decrypted)

	if decrypted == m {
		fmt.Println("SUCCESS!!!")
	} else {
		fmt.Println("something went wrong")
	}
}
