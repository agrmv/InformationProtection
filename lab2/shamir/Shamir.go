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

/*TODO THIS
type Shamir struct {
	ShamirKeys Keys
	p          int64
	m          int64
}*/

type Message struct {
	encryptMessage []int64
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

func generateAliceKeys(keys *Keys, p int64) {
	keys.AliceKeys = GenerateCD(&keys.AliceKeys, p) // генерируем взаимно-простые C и D для Алисы
}

func generateBobKeys(keys *Keys, p int64) {
	keys.BobKeys = GenerateCD(&keys.BobKeys, p) // генерируем взаимно-простые C и D для Боба
}

func Encrypt(alice Pair, m, p int64) int64 {
	return methods.ModularPow(methods.ModularPow(m, alice.First, p), alice.Second, p)
}

func Decrypt(bob Pair, p, encrypt int64) int64 {
	return methods.ModularPow(methods.ModularPow(encrypt, bob.First, p), bob.Second, p)
}

func EncryptMessage(file []byte, fileSize int64, keys Keys, message *Message, p int64) {
	message.encryptMessage = make([]int64, fileSize)
	for i, v := range file {
		message.encryptMessage[i] = Encrypt(keys.AliceKeys, int64(v), p)
	}
}

func DecryptMessage(keys Pair, message *Message, p int64) {
	message.decryptMessage = make([]byte, len(message.encryptMessage))
	for i, v := range message.encryptMessage {
		message.decryptMessage[i] = byte(Decrypt(keys, p, v))
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
	p := methods.DefaultGeneratePrime() // генерируем общее простое число
	keys := Keys{}
	generateAliceKeys(&keys, p)
	generateBobKeys(&keys, p)

	message := Message{}
	writeKeyToJson("lab2/shamir/resources/privateKeys.json", keys.BobKeys)
	file, fileSize := methods.ReadFile("lab2/resourcesGlobal/test.jpg")
	EncryptMessage(file, fileSize, keys, &message, p)
	DecryptMessage(getKeyFromJson("lab2/shamir/resources/privateKeys.json"), &message, p)
	methods.WriteFile("lab2/shamir/resources/decrypt.jpg", message.decryptMessage)
}
