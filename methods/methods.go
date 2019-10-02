package methods

import (
	"math/rand"
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

func Gcd(a, b int64) int64 {
	for b != 0 {
		r := a % b
		a, b = b, r
	}
	return a
}

func GcdExtended(a, b int64) (_, _, _ int64) {
	if a == 0 {
		return b, 0, 1
	}

	gcd, x, y := GcdExtended(b%a, a)
	return gcd, y - (b/a)*x, x
}

/*
	проверка на простоту
	n - значение для провреки
	к - число тестов на простоу
*/
func TestFerma(n int64, k int) bool {
	if n == 2 {
		return true
	}
	if n&1 != 1 {
		return false
	}
	for i := 0; i < 5; i++ {
		a := rand.Int63n(n-1) + 1
		if Gcd(a, n) != 1 || ModularPow(a, n-1, n) != 1 {
			return false
		}
	}
	return true
}
