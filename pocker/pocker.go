package main

import (
	"../methods"
	"fmt"
	"math/rand"
)

type Player struct {
	P           int64
	C           int64
	D           int64
	cardEncrypt []int64
	cardDecrypt []int64
}

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
	cards := make([]int64, 6)
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

func initPlayer(p int64) Player {
	c, d := GenerateKeys(p)
	return Player{P: p, C: c, D: d}
}

func initPlayers(p, n int64) []Player {
	var players []Player
	for i := int64(0); i < n; i++ {
		player := initPlayer(p)
		players = append(players, player)
		fmt.Println("player ", i, ": ", player)
	}
	return players
}

func (p *Player) encryptDeck(deck []int64) []int64 {
	encryptDeck := make([]int64, len(deck))
	for i, item := range deck {
		encryptDeck[i] = methods.ModularPow(item, p.C, p.P)
	}
	return encryptDeck
}

func (p *Player) decryptDeck(deck []int64) []int64 {
	decryptDeck := make([]int64, len(deck))
	for i, item := range deck {
		decryptDeck[i] = methods.ModularPow(item, p.D, p.P)
	}
	return decryptDeck
}

func encryptAllDeck(players []Player, encrypt []int64) []int64 {
	for _, player := range players {
		encrypt = player.encryptDeck(encrypt)
	}
	return encrypt
}

func decryptAllDeck(players []Player, decrypt []int64) []int64 {
	for _, player := range players {
		decrypt = player.decryptDeck(decrypt)
	}
	return decrypt
}

func distributeCards(players *[]Player, deck *[]int64) {
	//2 - число карт на игрока
	for i := 0; i < 2; i++ {
		for j, _ := range *players {
			if len(*deck)-1 < 0 {
				fmt.Println()
				break
			}
			(&(*players)[j]).cardEncrypt = append((&(*players)[j]).cardEncrypt, (*deck)[len(*deck)-1])
			*deck = (*deck)[:len(*deck)-1]
		}
	}
}

func distributeDecryptedCards(players *[]Player, deck *[]int64) {
	//2 - число карт на игрока
	for i := 0; i < 2; i++ {
		for j, _ := range *players {
			if len(*deck)-1 < 0 {
				fmt.Println()
				break
			}
			(&(*players)[j]).cardDecrypt = append((&(*players)[j]).cardDecrypt, (*deck)[len(*deck)-1])
			*deck = (*deck)[:len(*deck)-1]
		}
	}
}

func main() {
	P, _, _ := methods.GeneratePQg(53)

	//TODO FIX RANDOM
	arr := generateCards(P)
	fmt.Println("deck:", arr)

	players := initPlayers(P, 3)
	encryptDeck := encryptAllDeck(players, arr)
	fmt.Println("encrypt deck:", encryptDeck)

	decryptDeck := decryptAllDeck(players, encryptDeck)
	fmt.Println("decrypt deck:", decryptDeck)

	distributeCards(&players, &encryptDeck)
	fmt.Println("players with encrypted", players)

	//костыль
	distributeDecryptedCards(&players, &decryptDeck)
	fmt.Println("players with decrypted", players)

	//TODO DECRYPT
	//TODO change num to pic/text
}
