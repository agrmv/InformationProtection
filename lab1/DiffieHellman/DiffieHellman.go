package DiffieHellman

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
)

type Person struct {
	PrivateKey int64
	PublicKey  int64
	CommonKey  int64
}

func GeneratePQg(maxValue int64) (int64, int64, int64) {
	var Q = big.NewInt(1)
	var P = big.NewInt(1)
	var g = big.NewInt(1)
	var tmp = big.NewInt(1)
	for true {
		Q, err := rand.Int(rand.Reader, big.NewInt(maxValue/2-1))
		if err != nil {
			panic("ERROR in rand.Int")
		}
		P = P.Add(P.Mul(big.NewInt(2), Q), big.NewInt(1)) // P = 2 * Q + 1
		if Q.ProbablyPrime(20) && P.ProbablyPrime(20) {
			break
		}
	}
	for true {
		g, _ = rand.Int(rand.Reader, big.NewInt(P.Int64()-1))
		if g.Cmp(big.NewInt(1)) != 1 {
			continue
		}
		if tmp.Exp(g, Q, P) != big.NewInt(1) { // if g^Q mod P != 1
			break
		}
	}
	return P.Int64(), (P.Int64() - 1) / 2, g.Int64()
}

func GeneratePrivateKey(P int64) int64 {
	var key = big.NewInt(1)
	for true {
		key, _ = rand.Int(rand.Reader, big.NewInt(P))
		if key.Cmp(big.NewInt(1)) >= 0 {
			break
		}
	}
	return key.Int64()
}

func GeneratePublicKey(g int64, PrivateKey int64, P int64) int64 {
	var key = big.NewInt(1)
	bigG := big.NewInt(g)
	bigPrivateKey := big.NewInt(PrivateKey)
	bigP := big.NewInt(P)
	key = bigG.Exp(bigG, bigPrivateKey, bigP)

	return key.Int64()
}

func FindCommonKey(PublicKey int64, PrivateKey int64, P int64) int64 {
	Y := float64(PublicKey)
	X := float64(PrivateKey)

	commonKey := int64(math.Pow(Y, X)) % P
	return commonKey
}

func MainForDiffieHellman() {
	P, Q, g := GeneratePQg(10000)
	fmt.Printf("P = %d, Q = %d, g = %d\n", P, Q, g)

	var A, B Person
	A.PrivateKey = GeneratePrivateKey(P)
	A.PublicKey = GeneratePublicKey(g, A.PrivateKey, P)

	B.PrivateKey = GeneratePrivateKey(P)
	B.PublicKey = GeneratePublicKey(g, B.PrivateKey, P)

	A.CommonKey = FindCommonKey(B.PublicKey, A.PrivateKey, P)
	B.CommonKey = FindCommonKey(A.PublicKey, B.PrivateKey, P)

	fmt.Println("A: ", A)
	fmt.Println("B: ", B)
	if A.CommonKey == B.CommonKey {
		fmt.Printf("Ключи совпали: A = %d, B = %d", A.CommonKey, B.CommonKey)
	} else {
		fmt.Printf("Ключи НЕ совпали: A = %d, B = %d", A.CommonKey, B.CommonKey)
	}
}
