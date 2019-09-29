package main

import "fmt"

func isBigAlphabeticalLetter(letter int32) bool {
	if 65 <= letter && letter <= 90 {
		return true
	}
	return false
}

func isSmallAlphabeticalLetter(letter int32) bool {
	if 97 <= letter && letter <= 122 {
		return true
	}
	return false
}

func checkAlphabetical(in *string) {
	var out []rune
	for _, letter := range *in {
		if isBigAlphabeticalLetter(letter) {
			out = append(out, letter)
		} else if isSmallAlphabeticalLetter(letter) {
			out = append(out, letter-32)
		}
	}
	*in = string(out)
}

func encodePair(a, b rune) rune {
	return (((a - 'A') + (b - 'A')) % 26) + 'A'
}

func decodePair(a, b rune) rune {
	return ((((a - 'A') - (b - 'A')) + 26) % 26) + 'A'
}

func Encode(msg, key string) string {
	checkAlphabetical(&msg)
	checkAlphabetical(&key)
	out := make([]rune, 0, len(msg))
	for i, v := range msg {
		out = append(out, encodePair(v, rune(key[i%len(key)])))
	}
	return string(out)
}

func Decode(msg, key string) string {
	checkAlphabetical(&msg)
	checkAlphabetical(&key)
	out := make([]rune, 0, len(msg))
	for i, v := range msg {
		out = append(out, decodePair(v, rune(key[i%len(key)])))
	}
	return string(out)
}

func main() {
	key := "papich"
	message := "Vsem moim bratyam salam"

	encoded := Encode(message, key)
	decoded := Decode(encoded, key)

	fmt.Println(encoded)
	fmt.Println(decoded)
}
