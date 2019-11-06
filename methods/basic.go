package methods

import (
	"math"
)

func ModularPow(base, exp, module int64) int64 {
	result := int64(1)
	for ; exp > 0; exp >>= 1 {
		if exp&1 == 1 {
			result = (result * base) % module
		}
		base = (base * base) % module
	}
	return result
}

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

func Power(x int64, y int64) int64 {
	return int64(math.Pow(float64(x), float64(y)))
}
