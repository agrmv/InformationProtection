package Euclid

import (
	"fmt"
	"log"
)

func gcd(a int64, b int64, x *int64, y *int64) int64 {
	if a == 0 {
		*x = 0
		*y = 1
		return b
	}
	var x1, y1, d int64
	d = gcd(b%a, a, &x1, &y1)
	*x = y1 - (b/a)*x1
	*y = x1
	return d
}

func MainForEuclid() {
	var a, b, x, y int64
	fmt.Println("Input a, x, p")
	if _, err := fmt.Scan(&a, &b); err != nil {
		log.Print("  Scan for a, b failed ", err)
		return
	}
	fmt.Printf("\n\nResult: d = %d, x = %d, y = %d", gcd(a, b, &x, &y), x, y)
}
