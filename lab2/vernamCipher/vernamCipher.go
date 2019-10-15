package main

import (
	"../../methods"
	"encoding/json"
	"io/ioutil"
	"math/rand"
)

func Encrypt(message []byte, key string) []byte {
	out := make([]byte, 0, len(message))
	for i, v := range message {
		out = append(out, v^(key[i%len(key)]))
	}
	return out
}

func Decrypt(message []byte, key string) []byte {
	out := make([]byte, 0, len(message))
	for i, v := range message {
		out = append(out, v^(key[i%len(key)]))
	}
	return out
}

func generatePrivateKey(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func writeKeyToJson(path, key string) {
	privateKey, _ := json.Marshal(key)
	_ = ioutil.WriteFile(path, privateKey, 0644)
}

func getKeyFromJson(path string) string {
	file, _ := ioutil.ReadFile(path)
	var key string
	_ = json.Unmarshal(file, &key)
	return key
}

func main() {
	{
		file, _ := methods.ReadFile("lab2/resourcesGlobal/test.jpg")
		writeKeyToJson("lab2/vernamCipher/resources/privateKey.json", generatePrivateKey(10))
		encoded := Encrypt(file, getKeyFromJson("lab2/vernamCipher/resources/privateKey.json"))
		methods.WriteFile("lab2/vernamCipher/resources/encrypted.jpg", encoded)
	}
	{
		file, _ := methods.ReadFile("lab2/vernamCipher/resources/encrypted.jpg")
		decoded := Decrypt(file, getKeyFromJson("lab2/vernamCipher/resources/privateKey.json"))
		methods.WriteFile("lab2/vernamCipher/resources/decrypted.jpg", decoded)
	}
}
