package main

import (
	"../../methods"
	"fmt"
	rand2 "math/rand"
)

func GenerateCD(C *int64, D *int64, p int64) {
	*C = 1
	*D = rand2.Int63n(p)
	for *C < p {
		if (*C**D)%(p-1) == 1 {
			return
		}
		*C = *C + 1
	}
	if *C == p {
		GenerateCD(C, D, p)
	}
}

func Shamir(m int64) {
	var Ca, Da, p int64
	var Cb, Db int64

	p = methods.GeneratePrimes(100) // генерируем общее простое число

	GenerateCD(&Ca, &Da, p) // генерируем взаимно-простые C и D для Алисы
	fmt.Printf("\nCa = %d, Da = %d, p = %d", Ca, Da, p)

	GenerateCD(&Cb, &Db, p) // генерируем взаимно-простые C и D для Боба
	fmt.Printf("\nCb = %d, Db = %d, p = %d\n", Cb, Db, p)

	fmt.Printf("Secret message: %d\n", m)

	x1 := methods.ModularPow(m, Ca, p) //step 1
	fmt.Printf("\nstep1: %d\n", x1)

	x2 := methods.ModularPow(x1, Cb, p) //step 2
	fmt.Printf("step2: %d\n", x2)

	x3 := methods.ModularPow(x2, Da, p) //step 3
	fmt.Printf("step3: %d\n", x3)

	x4 := methods.ModularPow(x3, Db, p) //step 4
	fmt.Printf("step4: %d\n", x4)

	if x4 == m {
		fmt.Println("SUCCESS!!!")
	} else {
		fmt.Println("something went wrong")
	}
}
