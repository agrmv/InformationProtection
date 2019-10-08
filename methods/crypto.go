package methods

import (
	"math"
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

func GeneratePrime(generator func() int64) int64 {
	for {
		prime := generator()
		if TestFerma(prime, 3) {
			return prime
		}
	}
}

func DefaultGeneratePrime() int64 {
	return GeneratePrime(rand.Int63)
}
