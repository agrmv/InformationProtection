package euclid

import (
	"fmt"
	"log"
)

func Gcd(a int64, b int64) (result int64, x int64, y int64) {
	if a == 0 {
		return b, 0, 1
	}

	result, x, y = Gcd(b%a, a)
	return result, y - (b/a)*x, x
}

func MainForEuclid() {
	var a, b int64
	fmt.Println("Input a, x")
	if _, err := fmt.Scan(&a, &b); err != nil {
		log.Print("  Scan for a, b failed ", err)
		return
	}
	fmt.Println(Gcd(a, b))
}
