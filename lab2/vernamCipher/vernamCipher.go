package main

import "fmt"

func Encode(msg, key string) []rune {
	out := make([]rune, 0, len(msg))
	for i, v := range msg {
		out = append(out, v^rune(key[i%len(key)]))
	}
	return out
}

func Decode(msg []rune, key string) []rune {
	out := make([]rune, 0, len(msg))
	for i, v := range msg {
		out = append(out, v^rune(key[i%len(key)]))
	}
	return out
}

func main() {
	key := "papich"
	message := "Vsem moim bratyam salam"

	encoded := Encode(message, key)
	decoded := Decode(encoded, key)

	fmt.Println(encoded)
	fmt.Println(string(decoded))
}
