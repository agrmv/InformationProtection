package main

import (
	"../methods"
	"fmt"
	"math/rand"
)

func GenerateKeys(P int64) (int64, int64) {
	var c, d int64
	for {
		c = methods.LimitedGeneratePrime(53)
		_, d, _ = methods.GcdExtended(c, P-1)
		if d > 0 {
			break
		}
	}
	_, d, _ = methods.GcdExtended(c, P-1)
	return c, d
}

//func chooseRandomCard(a, b, c int64) int64 {
//	if rand.Int63n(10) % 3 == 0 {
//		return a
//	} else if rand.Int63n(10) % 3 == 1 {
//		return b
//	}
//	return c
//}

func chooseRandomCard(a []int64) int64 {
	return a[rand.Intn(len(a))]
}

func chooseRandomCard2(a, b int64) int64 {
	if rand.Int63n(10) == 0 {
		return a
	}
	return b
}

func find(slice []int64, element int64) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

func generateCards(P int64) []int64 {
	var r int64
	cards := make([]int64, 3)
	for i, _ := range cards {
		for {
			r = rand.Int63n(P - 1)
			if r >= 2 && !find(cards, r) {
				break
			}
		}
		cards[i] = r
	}
	return cards
}

func deleteCardFromDeck(deck []int64, card int64) []int64 {
	otherCards := make([]int64, 0)
	for index, _ := range deck {
		if deck[index] == card {
			continue
		}
		otherCards = append(otherCards, deck[index])
	}
	return otherCards
}

func encodeDeck(deck []int64, c, P int64) []int64 {
	encodeDeck := make([]int64, len(deck))
	for i, item := range deck {
		encodeDeck[i] = methods.ModularPow(item, c, P)
	}
	return encodeDeck
}

func main() {
	P, _, _ := methods.GeneratePQg(53)

	ca, da := GenerateKeys(P) // keys of Alice
	cb, db := GenerateKeys(P) // keys of Bob

	cards := generateCards(P) // our deck
	fmt.Print("cards: ")
	fmt.Print(cards)
	fmt.Print("\n")

	encodeCards := encodeDeck(cards, ca, P) // Alice encode the deck and send it to the Bob
	//fmt.Printf("u1, u2, u3: %d, %d, %d\n", encodeCards[0], encodeCards[1], encodeCards[2])

	B := chooseRandomCard(encodeCards) // Bob choose one card from deck and send it back to Alice
	//fmt.Printf("B: %d\n", B)

	Bdecode := methods.ModularPow(B, da, P) // Alice decode her card
	fmt.Printf("Alice's card: %d\n", Bdecode)

	otherCards := deleteCardFromDeck(encodeCards, B) // now Bob work with other cards
	//fmt.Println(otherCards)

	v1 := methods.ModularPow(otherCards[0], cb, P) // Bob generate numbers(decode decoding deck?) and send that to Alice
	v2 := methods.ModularPow(otherCards[1], cb, P)
	//fmt.Printf("v1, v2: %d, %d\n", v1, v2)

	A := chooseRandomCard2(v1, v2) // Alice choose one of them
	//fmt.Printf("A: %d\n", A)

	w1 := methods.ModularPow(A, da, P) // Alice decode chosen card and send that to Bob

	z := methods.ModularPow(w1, db, P) // Bob decode the card
	fmt.Printf("Bob's card: %d\n", z)
}
