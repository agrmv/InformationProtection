package main

import (
	"./babyGiantStep"
	"./diffieHellman"
	"./euclid"
	"./modularPow"
	"fmt"
)

func main() {
	fmt.Println("\n***modularPow***")
	modularPow.MainForModularPow()

	fmt.Println("\n\n***Euklid***")
	euclid.MainForEuclid()

	fmt.Println("\n\n***Diffie-Hellman***")
	diffieHellman.MainForDiffieHellman()

	fmt.Println("\n\n***BabyGiant***")
	babyGiantStep.MainForBabyGiantStep()
}
