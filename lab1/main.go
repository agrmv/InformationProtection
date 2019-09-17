package main

import (
	"./Euklid"
	"./modularPow"
	"fmt"
)
import "./DiffieHellman"

func main() {
	fmt.Println("\n***modularPow***")
	modularPow.MainForModularPow()

	fmt.Println("\n\n***Euklid***")
	Euklid.MainForEuklid()

	fmt.Println("\n\n***Diffie-Hellman***")
	DiffieHellman.MainForDiffieHellman()
}
