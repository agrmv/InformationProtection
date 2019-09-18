package main

import (
	"./DiffieHellman"
	"./Euklid"
	"./babyGiantStep"
	"./modularPow"
	"fmt"
)

func main() {
	fmt.Println("\n***modularPow***")
	modularPow.MainForModularPow()

	fmt.Println("\n\n***Euklid***")
	Euklid.MainForEuklid()

	fmt.Println("\n\n***Diffie-Hellman***")
	DiffieHellman.MainForDiffieHellman()

	fmt.Println("\n\n***BabyGiant***")
	babyGiantStep.MainForBabyGiantStep()
}
