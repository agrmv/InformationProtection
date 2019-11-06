package main

import (
	"crypto/rand"
	"crypto/sha1"
	"math/big"
)

type PublicKeys struct {
	P *big.Int
	Q *big.Int
	A *big.Int
}

type ExchangeData struct {
	Data PublicKeys
	Y    *big.Int
}

func generateA(b, q, p *big.Int) (a *big.Int) {
	for {
		g, err := rand.Int(rand.Reader, new(big.Int).Sub(p, big.NewInt(1)))
		if err != nil {
			panic(err)
		}
		a = new(big.Int).Exp(g, b, p)
		if a.Cmp(big.NewInt(1)) == 1 {
			return a
		}
	}
}

func generatePB(q *big.Int) (_, _ *big.Int) {
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(16), nil)
	for {
		b, _ := rand.Int(rand.Reader, max)
		p := b.Mul(b, q).Add(b, big.NewInt(1))
		if p.ProbablyPrime(50) {
			return p, b
		}
	}
}

func GenerateXY(d PublicKeys) (_private, _public *big.Int) {
	x, err := rand.Int(rand.Reader, d.Q)
	if err != nil {
		panic(err)
	}
	y := new(big.Int).Exp(d.A, x, d.P)
	return x, y
}

func InitPQA() PublicKeys {
	q, err := rand.Prime(rand.Reader, 15)
	if err != nil {
		panic(err)
	}
	p, b := generatePB(q)
	a := generateA(b, q, p)
	return PublicKeys{p, q, a}
}

func GenerateSignature(d PublicKeys, x *big.Int, message string) (_, _ *big.Int) {
	var r, s, h *big.Int
	h = new(big.Int).SetBytes(sha1.New().Sum([]byte(message)))
	for {
		k, err := rand.Int(rand.Reader, d.Q)
		if err != nil {
			panic(err)
		}
		r = new(big.Int).Exp(d.A, k, d.P)
		r.Mod(r, d.Q)

		if r.Cmp(big.NewInt(0)) == 0 {
			continue
		}

		kh := new(big.Int).Mul(k, h)
		xr := new(big.Int).Mul(x, r)
		s = new(big.Int).Add(kh, xr)
		s.Mod(s, d.Q)
		if s.Cmp(big.NewInt(0)) == 0 || d.Q.Cmp(r) != 1 && d.Q.Cmp(s) != 1 {
			continue
		}
		return r, s
	}
}

func CheckSignatureGost(r, s *big.Int, m string, d ExchangeData) bool {
	h := new(big.Int).SetBytes(sha1.New().Sum([]byte(m)))
	invers_h := new(big.Int).ModInverse(h, d.Data.Q) //h^-1
	println("h^-1 = " + invers_h.Text(10))

	sh := new(big.Int).Mul(s, invers_h)
	rh := new(big.Int).Mul(r, invers_h)
	rh.Neg(rh)
	u1 := new(big.Int).Mod(sh, d.Data.Q)
	u2 := new(big.Int).Mod(rh, d.Data.Q)
	au1 := new(big.Int).Exp(d.Data.A, u1, d.Data.P)
	yu2 := new(big.Int).Exp(d.Y, u2, d.Data.P)
	v := new(big.Int).Mul(au1, yu2)
	v.Mod(v, d.Data.P)
	v.Mod(v, d.Data.Q)
	return v.Cmp(r) != 0
}

func SignatureGostTEST(d PublicKeys, x *big.Int) (r, s *big.Int) {
	var k, h *big.Int
	h = big.NewInt(4)
	for {
		k, _ = rand.Int(rand.Reader, d.Q)
		r = new(big.Int).Exp(d.A, k, d.P)
		r.Mod(r, d.Q)
		if r.Cmp(big.NewInt(0)) == 0 {
			continue
		}

		kh := new(big.Int).Mul(k, h)
		xr := new(big.Int).Mul(x, r)
		s = new(big.Int).Add(kh, xr)
		s.Mod(s, d.Q)
		if s.Cmp(big.NewInt(0)) == 0 {
			continue
		}
		if d.Q.Cmp(r) != 1 && d.Q.Cmp(s) != 1 {
			continue
		}
		return r, s
	}
}

func ChecSignatureGostTEST(r, s *big.Int, d ExchangeData) bool {
	h := big.NewInt(4)
	invers_h := new(big.Int).ModInverse(h, d.Data.Q) //??????????
	println(invers_h.Text(10))

	sh := new(big.Int).Mul(s, invers_h)
	rh := new(big.Int).Mul(r, invers_h)
	rh.Neg(rh)
	u1 := new(big.Int).Mod(sh, d.Data.Q)
	u2 := new(big.Int).Mod(rh, d.Data.Q)
	au1 := new(big.Int).Exp(d.Data.A, u1, d.Data.P)
	yu2 := new(big.Int).Exp(d.Y, u2, d.Data.P)
	v := new(big.Int).Mul(au1, yu2)
	v.Mod(v, d.Data.P)
	v.Mod(v, d.Data.Q)
	return v.Cmp(r) == 0
}

func main() {

	/*	//for {
		publicKeys := InitPQA()

		println(" ======= \n", publicKeys.Q.Text(10), "   ", publicKeys.Q.BitLen())
		println(publicKeys.P.Text(10), "   ", publicKeys.P.BitLen())
		println(publicKeys.A.Text(10), "   ", publicKeys.P.BitLen(), " \n ======= ")
		//2994
		//42
		pubA, pivA := GenerateXY(publicKeys)
		println("a = " + pubA.Text(10))
		println(pivA.Text(10))
		r, s := SignatureGostTEST(publicKeys, pivA)
		println("r = " + r.Text(10))
		println("s = " + s.Text(10))
		lol := ExchangeData{publicKeys, pubA}
		check := ChecSignatureGostTEST(r, s, lol)
		println(check)
	//}*/

	PublicKeys := InitPQA()
	println("p = " + PublicKeys.P.Text(10))
	println("q = " + PublicKeys.Q.Text(10))
	println("a = " + PublicKeys.A.Text(10))

	// Alice
	X, Y := GenerateXY(PublicKeys)
	println("X: " + X.Text(10))
	println("Y: " + Y.Text(10))

	data := ExchangeData{PublicKeys, Y}
	m := "hello"
	r, s := GenerateSignature(PublicKeys, X, m)
	println(s)
	// Bob
	check := CheckSignatureGost(r, s, m, data)
	println(check)

}
