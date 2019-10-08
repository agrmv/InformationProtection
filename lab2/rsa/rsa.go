//package rsa
package main

import (
	"../../methods"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
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
	generated := rand.Int63() % 1000
	for !methods.TestFerma(generated, 5) {
		generated = rand.Int63() % 1000
	}
	return generated
}

//(1 < e < phi)
func getOpenExp(phi int64) int64 {
	generated := rand.Int63n(phi-1) + 1
	for !methods.TestFerma(generated, 5) {
		generated = rand.Int63n(phi-1) + 1
	}
	return generated
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

func writeKeyToJson(path string, keys Keys) {
	privateKeys, _ := json.Marshal(keys.PrivateKey)
	_ = ioutil.WriteFile(path, privateKeys, 0644)
}

func getKeyFromJson(path string) Pair {
	file, _ := ioutil.ReadFile(path)
	keys := Pair{}
	_ = json.Unmarshal(file, &keys)
	return keys
}

func EncryptMessage(file []byte, fileSize int64, keys Keys, message *Message) {
	message.encryptMessage = make([]int64, fileSize)
	for i, v := range file {
		message.encryptMessage[i] = Encrypt(keys.PublicKey, int64(v))
	}
}

func DecryptMessage(keys Pair, message *Message) {
	message.decryptMessage = make([]byte, len(message.encryptMessage))
	for i, v := range message.encryptMessage {
		message.decryptMessage[i] = byte(Decrypt(keys, v))
	}
}

func main() {

	rand.Seed(time.Now().UnixNano())
	message := Message{}

	fmt.Print("Choose n option:\n1)Encrypt\n2)Decrypt\n:")
	var option int
	_, _ = fmt.Fscan(os.Stdin, &option)
	switch option {
	case 1:
		keys := generateKeys()
		writeKeyToJson("lab2/rsa/resources/privateKeys.json", keys)
		file, fileSize := methods.ReadFile("lab2/resourcesGlobal/test.jpg")
		EncryptMessage(file, fileSize, keys, &message)
		DecryptMessage(getKeyFromJson("lab2/rsa/resources/privateKeys.json"), &message)
		methods.WriteFile("lab2/rsa/resources/decrypt.jpg", message.decryptMessage)
	case 2:
		/*file, fileSize := methods.ReadFile("lab2/rsa/resources/encrypt.jpg")
		DecryptMessage(file, fileSize, getKeyFromJson("lab2/rsa/resources/privateKeys.json"), &message)
		methods.WriteFile("lab2/rsa/resources/decrypt.jpg", message.decryptMessage)*/
	default:
		fmt.Println("Incorrect option")
	}
}
