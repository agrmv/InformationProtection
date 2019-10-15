package main

import (
	"../../lab1/diffieHellman"
	"../../methods"
	"fmt"
	"math"
	"math/rand"
)

type Alice struct {
	message int64
	k       int64
	a       int64
	b       int64
}

type Bob struct {
	decodedMessage int64
	x              int64
	y              int64
}

func generateSessionKey(P int64) int64 {
	var key int64
	for true {
		key = rand.Int63n(P)

		if key > 1 {
			break
		}
	}
	return key
}

func main() {
	var A Alice
	A.message = 10 //strconv.ParseInt(os.Args[1], 10, 64)
	var P, g int64
	for true {
		P, _, g = diffieHellman.GeneratePQg(20)
		if P > A.message {
			break
		}
	}
	var B Bob
	B.x = generateSessionKey(P - 1)
	B.y = methods.ModularPow(g, B.x, P)

	A.k = generateSessionKey(P)

	A.a = methods.ModularPow(g, A.k, P)
	yk := int64(math.Pow(float64(B.y), float64(A.k)))
	A.b = (yk * A.message) % P

	fmt.Printf("A.k = %d, A.a = %d, A.b = %d\n", A.k, A.a, A.b)
	fmt.Printf("B.x = %d, B.y = %d, P = %d\n", B.x, B.y, P)

	aa := int64(math.Pow(float64(A.a), float64(P-1-B.x)))
	B.decodedMessage = (A.b * aa) % P

	fmt.Printf("Message of Alice: %d\nDecoded message of Bob: %d\n", A.message, B.decodedMessage)
}
