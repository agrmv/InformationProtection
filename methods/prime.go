package methods

import (
	"math/rand"
)

func generatePrimeDefault(generator func() int64) int64 {
	for {
		prime := generator()
		if TestFerma(prime, 5) {
			return prime
		}
	}
}

func generatePrimeLimited(limit int64) int64 {
	for {
		prime := rand.Int63n(limit)
		if prime == 1 {
			continue
		}
		if TestFerma(prime, 20) {
			return prime
		}
	}
}

func DefaultGeneratePrime() int64 {
	return generatePrimeDefault(rand.Int63)
}

func LimitedGeneratePrime(max int64) int64 {
	return generatePrimeLimited(max)
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
