package main

import (
	"./DiffieHellman"
	"./Euclid"
	"./babyGiantStep"
	"fmt"
)

func main() {
	/*fmt.Println("\n***modularPow***")
	modularPow.MainForModularPow()*/

	fmt.Println("\n\n***Euklid***")
	Euclid.MainForEuclid()

	fmt.Println("\n\n***Diffie-Hellman***")
	DiffieHellman.MainForDiffieHellman()

	fmt.Println("\n\n***BabyGiant***")
	babyGiantStep.MainForBabyGiantStep()
}
