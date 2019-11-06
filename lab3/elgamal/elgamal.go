package main

import (
	"../../methods"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
)

const (
	text  = "resourcesGlobal/test.txt"
	image = "resourcesGlobal/test.jpg"
)

func generateS(k, u, P int64) int64 {
	_, antiK, _ := methods.GcdExtended(k, P-1)
	return methods.ModularPow(antiK*u, 1, P-1)
}

func writeSignature(s int64) {
	f, err := os.Create("lab3/elgamal/sign.txt")
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.WriteString(fmt.Sprintf("%d", s))
	if err != nil {
		log.Fatal(err)
	} // write signature in file
}

func encode(P, g, secretKey int64) {

}

func matchSignature(y, r, g, P int64) bool {
	file, _ := methods.ReadFile(text)

	h := sha1.New()
	h.Write(file)
	sha1_hash := binary.BigEndian.Uint64(h.Sum(nil))
	fmt.Println(sha1_hash)
	v1 := power(y, r) * power(r, int64(sha1_hash))
	v2 := methods.ModularPow(g, int64(sha1_hash), P)
	return v1 == v2
}

func main() {
	// generate public params
	P, _, g := methods.GeneratePQg(256)
	secretKey := rand.Int63n(P - 2)
	y := methods.ModularPow(g, secretKey, P)

	file, _ := methods.ReadFile(text)

	// encode
	h := sha1.New()
	h.Write(file)
	sha1_hash := binary.BigEndian.Uint64(h.Sum(nil))
	fmt.Println(sha1_hash)
	// генерируем простое число меньше P
	// оно автоматически удовлетворяет условию НОД(P,k) = 1
	k := methods.LimitedGeneratePrime(P - 1)
	r := methods.ModularPow(g, k, P)

	//u := methods.ModularPow(int64(sha1_hash) - secretKey * r, 1, P - 1)
	u := (int64(sha1_hash) - secretKey*r) % (P - 1)
	s := generateS(k, u, P)

	writeSignature(s)

	//decode
	fmt.Println(matchSignature(y, r, g, P))
}
