package methods

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
