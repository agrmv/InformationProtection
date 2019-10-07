package main

import (
	"../../methods"
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

	file, _ := methods.ReadFile("lab2/resourcesGlobal/test.jpg")
	key, _ := methods.ReadFile("lab2/vernamCipher/resources/key.txt")

	encoded := Encode(file, string(key))
	decoded := Decode(encoded, string(key))

	methods.WriteFile("lab2/vernamCipher/resources/encode.jpg", encoded)
	methods.WriteFile("lab2/vernamCipher/resources/decode.jpg", decoded)
}
