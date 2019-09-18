package babyGiantStep

import (
	"fmt"
	"math"
)
import . "../modularPow"

func BabyGiantStep(a int64, y int64, p int64) int64 { // a^x = b % p
	n := int64(math.Sqrt(float64(p))) + 1
	var values = make(map[int64]int64)

	for i := n; i >= 1; i-- {
		values[ModularPow(a, i*n, p)] = i
	}

	for i := int64(0); i <= n; i++ {
		cur := (ModularPow(a, i, p) * y) % p
		if val, ok := values[cur]; ok {
			answer := val*n - i
			if answer < p {
				return answer
			}
		}
	}
	return -1
}

func MainForBabyGiantStep() {
	a := int64(7)
	p := int64(8)
	y := int64(7)

	fmt.Println(BabyGiantStep(a, y, p))

}
