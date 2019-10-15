//package rsa
package main

import (
	"../../methods"
	"encoding/json"
	"io/ioutil"
)

type Pair struct {
	First  int64 `json:"First"`
	Second int64 `json:"Second"`
}

type Keys struct {
	PublicKey  Pair `json:"PublicKey"`
	PrivateKey Pair `json:"PrivateKey"`
}

type Message struct {
	encryptMessage []int64
	decryptMessage []byte
}

func generatePrivateKey() int64 {
	return methods.LimitedGeneratePrime(1000)
}

//(1 < e < phi)
func getOpenExp(phi int64) int64 {
	return methods.LimitedGeneratePrime(phi)
}

//d = e^-1 mod phi
func getSecretExp(e, phi int64) int64 {
	_, inverse, _ := methods.GcdExtended(e, phi)
	for inverse < 0 {
		inverse += phi
	}
	return inverse
}

//https://ru.wikipedia.org/wiki/RSA
func generateKeys() Keys {
	var result Keys
	var p, q int64

	p = generatePrivateKey()
	q = generatePrivateKey()

	n := p * q
	// вычисление функции эйлера
	phi := (p - 1) * (q - 1)
	e := getOpenExp(phi)

	result.PublicKey = Pair{e, n}

	//вычисление секретной экспаненты
	d := getSecretExp(e, phi)

	result.PrivateKey = Pair{d, n}
	return result
}

func Encrypt(key Pair, value int64) int64 {
	return methods.ModularPow(value, key.First, key.Second)
}

func Decrypt(key Pair, value int64) int64 {
	return methods.ModularPow(value, key.First, key.Second)
}

func writeKeyToJson(path string, keys Pair) {
	privateKeys, _ := json.Marshal(keys)
	_ = ioutil.WriteFile(path, privateKeys, 0644)
}

func getKeyFromJson(path string) Pair {
	file, _ := ioutil.ReadFile(path)
	keys := Pair{}
	_ = json.Unmarshal(file, &keys)
	return keys
}

func EncryptMessage(file []byte, fileSize int64, keys Pair, message *Message) {
	message.encryptMessage = make([]int64, fileSize)
	for i, v := range file {
		message.encryptMessage[i] = Encrypt(keys, int64(v))
	}
}

func DecryptMessage(keys Pair, message *Message) {
	message.decryptMessage = make([]byte, len(message.encryptMessage))
	for i, v := range message.encryptMessage {
		message.decryptMessage[i] = byte(Decrypt(keys, v))
	}
}

func main() {

	keys := generateKeys()
	message := Message{}
	{
		writeKeyToJson("lab2/rsa/resources/privateKeys.json", keys.PrivateKey)
		file, fileSize := methods.ReadFile("lab2/resourcesGlobal/test.jpg")
		EncryptMessage(file, fileSize, keys.PublicKey, &message)
		//methods.WriteFile("lab2/rsa/resources/encrypt.jpg", message.encryptMessage)
	}
	{
		//file, fileSize := methods.ReadFile("lab2/rsa/resources/encrypt.jpg")
		DecryptMessage(getKeyFromJson("lab2/rsa/resources/privateKeys.json"), &message)
		methods.WriteFile("lab2/rsa/resources/decrypt.jpg", message.decryptMessage)
	}
}
