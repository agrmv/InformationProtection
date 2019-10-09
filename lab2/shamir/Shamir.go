package main

import (
	"../../methods"
	"encoding/json"
	"io/ioutil"
	rand2 "math/rand"
)

type Pair struct {
	First  int64 `json:"First"`
	Second int64 `json:"Second"`
}

type Keys struct {
	AliceKeys Pair `json:"PublicKey"`
	BobKeys   Pair `json:"PrivateKey"`
}

type Message struct {
	encryptMessage []byte
	decryptMessage []byte
}

func GenerateCD(keys *Pair, p int64) Pair {
	keys.First = int64(1)
	keys.Second = rand2.Int63n(p)
	for keys.First < p {
		if (keys.First*keys.Second)%(p-1) == 1 {
			return Pair{keys.First, keys.Second}
		}
		keys.First += 1
	}
	if keys.First == p {
		GenerateCD(keys, p)
	}
	return Pair{keys.First, keys.Second}
}

func generateKeys(keys *Keys, p int64) Keys {
	return Keys{GenerateCD(&keys.AliceKeys, p), GenerateCD(&keys.BobKeys, p)}
}

func Encrypt(alice Pair, m, p int64) int64 {
	return methods.ModularPow(methods.ModularPow(m, alice.First, p), alice.Second, p)
}

func Decrypt(bob Pair, p, encrypt int64) int64 {
	return methods.ModularPow(methods.ModularPow(encrypt, bob.First, p), bob.Second, p)
}

func EncryptMessage(file []byte, fileSize int64, keys Pair, message *Message, p int64) {
	message.encryptMessage = make([]byte, fileSize)
	for i, v := range file {
		message.encryptMessage[i] = byte(Encrypt(keys, int64(v), p))
	}
}

func DecryptMessage(file []byte, fileSize int64, keys Pair, message *Message, p int64) {
	message.decryptMessage = make([]byte, fileSize)
	for i, v := range file {
		message.decryptMessage[i] = byte(Decrypt(keys, p, int64(v)))
	}
}

func writeKeyToJson(path string, bob Pair) {
	privateKeys, _ := json.Marshal(bob)
	_ = ioutil.WriteFile(path, privateKeys, 0644)
}

func getKeyFromJson(path string) Pair {
	file, _ := ioutil.ReadFile(path)
	keys := Pair{}
	_ = json.Unmarshal(file, &keys)
	return keys
}

func main() {
	p := methods.LimitedGeneratePrime(1000) // генерируем общее простое число
	keys := generateKeys(&Keys{}, p)
	message := Message{}
	{
		writeKeyToJson("lab2/shamir/resources/privateKeys.json", keys.BobKeys)
		file, fileSize := methods.ReadFile("lab2/resourcesGlobal/test.jpg")
		EncryptMessage(file, fileSize, keys.AliceKeys, &message, p)
		methods.WriteFile("lab2/shamir/resources/encrypt.jpg", message.encryptMessage)
	}
	{
		file, fileSize := methods.ReadFile("lab2/shamir/resources/encrypt.jpg")
		DecryptMessage(file, fileSize, getKeyFromJson("lab2/shamir/resources/privateKeys.json"), &message, p)
		methods.WriteFile("lab2/shamir/resources/decrypt.jpg", message.decryptMessage)
	}
}
