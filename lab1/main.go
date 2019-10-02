package main

import (
	"./babyGiantStep"
	"./diffieHellman"
	"fmt"
)

func main() {
	fmt.Println("\n\n***Diffie-Hellman***")
	diffieHellman.MainForDiffieHellman()

	fmt.Println("\n\n***BabyGiant***")
	babyGiantStep.MainForBabyGiantStep()
}
