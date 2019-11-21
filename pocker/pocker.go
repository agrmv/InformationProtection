package main

import (
	"../methods"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

//const (
//	piki_2 = iota + 1
//	piki_3
//	piki_4
//	piki_5
//	piki_6
//	piki_7
//	piki_8
//	piki_9
//	piki_10
//	piki_v
//	piki_d
//	piki_k
//	piki_t
//
//	kresti_2
//	kresti_3
//	kresti_4
//	kresti_5
//	kresti_6
//	kresti_7
//	kresti_8
//	kresti_9
//	kresti_10
//	kresti_v
//	kresti_d
//	kresti_k
//	kresti_t
//
//	bubi_2
//	bubi_3
//	bubi_4
//	bubi_5
//	bubi_6
//	bubi_7
//	bubi_8
//	bubi_9
//	bubi_10
//	bubi_v
//	bubi_d
//	bubi_k
//	bubi_t
//
//)

var suit = []string{"♣", "♠", "♥", "♦"}
var high = []string{"J", "Q", "K", "A"}

func InitDeck(poker bool) (cart_deck []string) {
	var intiNum int
	if poker {
		intiNum = 2
	} else {
		intiNum = 6
	}

	for _, v := range suit {
		for i := intiNum; i <= 10; i++ {
			cart := strconv.FormatInt(int64(i), 10) + v
			cart_deck = append(cart_deck, cart)
		}
		for _, h := range high {
			cart := h + v
			cart_deck = append(cart_deck, cart)
		}
	}
	return
}

type Player struct {
	P           int64
	C           int64
	D           int64
	cardEncrypt []int64
	CardDecrypt []int64
}

func GenerateKeys(P int64) (int64, int64) {
	var c, d int64
	for {
		c = methods.LimitedGeneratePrime(P)
		_, d, _ = methods.GcdExtended(c, P-1)
		if d > 2 {
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
	cardsAll := make([]int64, 46)
	for i, _ := range cardsAll {
		cardsAll[i] = int64(i + 1)
	}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	r.Shuffle(len(cardsAll), func(i, j int) {

		cardsAll[i], cardsAll[j] = cardsAll[j], cardsAll[i]
	})
	return cardsAll
}

func initPlayer(p int64) Player {
	c, d := GenerateKeys(p)
	return Player{P: p, C: c, D: d}
}

func initPlayers(p, n int64) []Player {
	if n > 23 {
		panic("wrong num of players")
	}
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
			(&(*players)[j]).CardDecrypt = append((&(*players)[j]).CardDecrypt, (*deck)[len(*deck)-1])
			//fmt.Println(changeNumToString((&(*players)[j]).cardDecrypt))
			*deck = (*deck)[:len(*deck)-1]
		}
	}
}

func changeNumToString(arr []int64) []string {
	str := make([]string, len(arr))
	cards := InitDeck(true)
	for i, j := range arr {
		str[i] = cards[j-1]
	}
	return str
}

func main() {
	var P int64
	for {
		P, _, _ = methods.GeneratePQg(500)
		if P > 53 {
			break
		}
	}

	//TODO FIX RANDOM
	arr := generateCards(P)
	fmt.Println("deck:", changeNumToString(arr))

	players := initPlayers(P, 23)
	encryptDeck := encryptAllDeck(players, arr)
	fmt.Println("encrypt deck:", encryptDeck)

	decryptDeck := decryptAllDeck(players, encryptDeck)
	fmt.Println("decrypt deck:", changeNumToString(decryptDeck))

	distributeCards(&players, &encryptDeck)
	fmt.Println("players with encrypted", players)

	//костыль
	distributeDecryptedCards(&players, &decryptDeck)
	fmt.Println("players with decrypted", players)

	//TODO change num to pic/text

	for _, item := range players {
		fmt.Println(changeNumToString(item.CardDecrypt))
	}

}
