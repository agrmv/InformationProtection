package main

import (
	"../../methods"
	"crypto/sha1"
	"fmt"
	"math"
	"math/big"
	"math/rand"
)

const b = 4

type Person struct {
	secretKey int64
	publicKey int64
}

func generatePQ() (int64, int64) {
	var p, q int64
	for {
		q = methods.LimitedGeneratePrime(math.MaxInt32)
		if q < 255 {
			continue
		}
		p = b*q + 1
		if methods.TestFerma(p, 10) {
			break
		}
	}
	return p, q
}
func generateA(p int64) int64 {
	var a, g int64
	for {
		g = rand.Int63n(p)
		a = methods.ModularPow(g, b, p)
		if a > 1 {
			return a
		}

	}
}

func main() {
	//init public keys
	p, q := generatePQ() // common params
	a := generateA(p)    // common param

	Alice := new(Person)
	Bob := new(Person)

	Alice.secretKey = rand.Int63n(q)
	Alice.publicKey = methods.ModularPow(a, Alice.secretKey, p)

	Bob.secretKey = rand.Int63n(q)
	Bob.publicKey = methods.ModularPow(a, Bob.secretKey, p)

	file, _ := methods.ReadFile("resourcesGlobal/test.txt")
	/* ENCODE START */
	hash := sha1.New()
	hash.Write(file)
	sha1_hash := hash.Sum(nil)
	fmt.Println(sha1_hash[0])

	//generate signature
	//TODO move in func
	var r, s int64
	for {
		k := rand.Int63n(q)
		r = methods.ModularPow(a, k, p) % q
		if r == 0 {
			continue
		}
		s = (k*int64(sha1_hash[0]) + Alice.secretKey*r) % q

		if s != 0 {
			break
		}
	}
	/* ENCODE END*/

	/* DECODE START */
	hash2 := sha1.New()
	hash2.Write(file)
	sha1_hash2 := hash2.Sum(nil)
	fmt.Println(sha1_hash2[0])

	if r >= 0 && s >= q {
		panic("wrong r or s")
	}
	_, antiH, _ := methods.GcdExtended(int64(sha1_hash2[0]), q)
	fmt.Println((int64(sha1_hash2[0]) * antiH) % q)
	//u1 := s * antiH
	//u2 := (-r * antiH) % q

	//v := (methods.Power(a, u1) * methods.Power(Bob.publicKey, u2)) % p % q
	//v1 := new(big.Int).Exp(big.NewInt(a), big.NewInt(u1), nil)

	u1 := new(big.Int).Mod(new(big.Int).Mul(big.NewInt(s), big.NewInt(antiH)), big.NewInt(q))
	u2 := new(big.Int).Mod(new(big.Int).Mul(new(big.Int).Neg(big.NewInt(r)), big.NewInt(antiH)), big.NewInt(q))
	au1 := new(big.Int).Exp(big.NewInt(a), u1, big.NewInt(p))
	yu2 := new(big.Int).Exp(big.NewInt(Alice.publicKey), u2, big.NewInt(p))
	v := new(big.Int).Mod(new(big.Int).Mod(new(big.Int).Mul(au1, yu2), big.NewInt(p)), big.NewInt(q))
	//must be true
	fmt.Println(r)
	fmt.Println(v)

	/* DECODE END */
}
