package main

import (
	"../../methods"
	"encoding/json"
	"fmt"
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

func generateKeys(keys *Keys, p int64) Keys {
	return Keys{GenerateCD(&keys.AliceKeys, p), GenerateCD(&keys.BobKeys, p)}
}

func Encrypt(keys Keys, m, p int64) int64 {
	return methods.ModularPow(methods.ModularPow(m, keys.AliceKeys.First, p), keys.BobKeys.First, p)
}

func Decrypt(keys Keys, encrypt, p int64) int64 {
	return methods.ModularPow(methods.ModularPow(encrypt, keys.AliceKeys.Second, p), keys.BobKeys.Second, p)
}

func EncryptMessage(file []byte, fileSize int64, keys Keys, message *Message, p int64) {
	message.encryptMessage = make([]int64, fileSize)
	for i, v := range file {
		message.encryptMessage[i] = Encrypt(keys, int64(v), p)
	}
	fmt.Println(message.encryptMessage)
}

func DecryptMessage(keys Keys, message *Message, p int64) {
	message.decryptMessage = make([]byte, len(message.encryptMessage))
	for i, v := range message.encryptMessage {
		message.decryptMessage[i] = byte(Decrypt(keys, v, p))
	}
	fmt.Println(message.decryptMessage)
}

func writeKeyToJson(path string, keys Keys) {
	privateKeys, _ := json.Marshal(keys)
	_ = ioutil.WriteFile(path, privateKeys, 0644)
}

func getKeyFromJson(path string) Keys {
	file, _ := ioutil.ReadFile(path)
	keys := Keys{}
	_ = json.Unmarshal(file, &keys)
	return keys
}

func main() {
	p := methods.LimitedGeneratePrime(1000) // генерируем общее простое число
	keys := generateKeys(&Keys{}, p)

	message := Message{}
	{
		writeKeyToJson("lab2/shamir/resources/privateKeys.json", keys)
		file, fileSize := methods.ReadFile("lab2/resourcesGlobal/test.jpg")
		EncryptMessage(file, fileSize, keys, &message, p)
		//methods.WriteFile("lab2/shamir/resources/encrypt.txt", message.encryptMessage)
	}
	{
		//file, fileSize := methods.ReadFile("lab2/shamir/resources/encrypt.txt")
		DecryptMessage(getKeyFromJson("lab2/shamir/resources/privateKeys.json"), &message, p)
		methods.WriteFile("lab2/shamir/resources/decrypt.jpg", message.decryptMessage)
	}
}
