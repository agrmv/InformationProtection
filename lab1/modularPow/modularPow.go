package main

import (
	"fmt"
	"log"
)

func modularPow(base *int64, exp *int64, modulus *int64) int64 {
	var result int64 = 1
	for *exp != 0 {
		if (*exp & 1) != 0 {
			result = (result * *base) % *modulus
		}
		*base = (*base * *base) % *modulus
		*exp >>= 1 //эквивалентно делению на 2^N
	}
	return result
}

func main() {
	fmt.Println("a^x mod p")
	fmt.Println("Input a, x, p")
	var a, x, p int64
	if _, err := fmt.Scan(&a, &x, &p); err != nil {
		log.Print("  Scan for a, x, p failed ", err)
		return
	}
	fmt.Printf("\n\nResult: %d", modularPow(&a, &x, &p))
}
