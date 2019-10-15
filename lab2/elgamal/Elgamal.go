package main

import (
	"../../methods"
	"fmt"
	"math"
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

func main() {
	var A Alice
	A.message = 10 //strconv.ParseInt(os.Args[1], 10, 64)
	var P, g int64
	for true {
		P, _, g = methods.GeneratePQg(20)
		if P > A.message {
			break
		}
	}
	var B Bob
	B.x = methods.LimitedGeneratePrime(P)
	B.y = methods.ModularPow(g, B.x, P)

	A.k = methods.LimitedGeneratePrime(P)

	A.a = methods.ModularPow(g, A.k, P)
	yk := int64(math.Pow(float64(B.y), float64(A.k)))
	A.b = (yk * A.message) % P

	fmt.Printf("A.k = %d, A.a = %d, A.b = %d\n", A.k, A.a, A.b)
	fmt.Printf("B.x = %d, B.y = %d, P = %d\n", B.x, B.y, P)

	aa := int64(math.Pow(float64(A.a), float64(P-1-B.x)))
	B.decodedMessage = (A.b * aa) % P

	fmt.Printf("Message of Alice: %d\nDecoded message of Bob: %d\n", A.message, B.decodedMessage)
}
