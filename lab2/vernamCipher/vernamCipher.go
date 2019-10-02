package main

import (
	"../renameMe"
)

func Encode(message []byte, key string) []byte {
	out := make([]byte, 0, len(message))
	for i, v := range message {
		out = append(out, v^(key[i%len(key)]))
	}
	return out
}

func Decode(message []byte, key string) []byte {
	out := make([]byte, 0, len(message))
	for i, v := range message {
		out = append(out, v^(key[i%len(key)]))
	}
	return out
}

func main() {

	file, _ := renameMe.ReadFile("lab2/resourcesGlobal/test.txt")
	key, _ := renameMe.ReadFile("lab2/vernamCipher/resources/key.txt")

	encoded := Encode(file, string(key))
	decoded := Decode(encoded, string(key))

	renameMe.WriteFile(decoded)
}
