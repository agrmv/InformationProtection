package main

import (
	"../../methods"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
)

type Pair struct {
	First  int64 `json:"First"`
	Second int64 `json:"Second"`
}

type Keys struct {
	PublicKey  Pair `json:"PublicKey"`
	PrivateKey Pair `json:"PrivateKey"`
}

type Alice struct {
	message int64
	k       int64
	a       int64
	b       int64
}

type Bob struct {
	decodedMessage int64
	x              int64
	y              int64
}

func generateAliceKeys(alice *Alice, P int64, g int64) {
	alice.k = methods.LimitedGeneratePrime(P)
	alice.a = methods.ModularPow(g, alice.k, P)
}

func generateBobKeys(bob *Bob, P int64, g int64) {
	bob.x = methods.LimitedGeneratePrime(P)
	bob.y = methods.ModularPow(g, bob.x, P)
}

func encrypt(alice *Alice, bob *Bob, P int64) {
	alice.b = (int64(math.Pow(float64(bob.y), float64(alice.k))) * alice.message) % P
}

/*func EncryptMessage(file []byte, fileSize int64, alice *Alice, bob *Bob, P int64) {
	alice.message = make([]int64, fileSize)
	for i, v := range file {
		alice.message[i] = encrypt(alice, bob, P, int64(v))
	}
}*/

func decrypt(alice *Alice, bob *Bob, P int64) {
	bob.decodedMessage = (alice.b * int64(math.Pow(float64(alice.a), float64(P-1-bob.x)))) % P
}

/*func DecryptMessage(alice *Alice, bob *Bob, P int64) {
	bob.decodedMessage = make([]byte, len(alice.message))
	for i, v := range alice.message {
		bob.decodedMessage[i] = byte(decrypt(alice, bob, P, v))
	}
}*/

func writeKeyToJson(path string, keys Keys) {
	marshalKeys, _ := json.Marshal(keys)
	_ = ioutil.WriteFile(path, marshalKeys, 0644)
}

func getKeyFromJson(path string) Keys {
	file, _ := ioutil.ReadFile(path)
	keys := Keys{}
	_ = json.Unmarshal(file, &keys)
	return keys
}

func main() {
	var A Alice
	var B Bob
	var P, g int64
	A.message = 11

	for true {
		P, _, g = methods.GeneratePQg(29)
		if P > A.message {
			break
		}
	}

	generateAliceKeys(&A, P, g)
	generateBobKeys(&B, P, g)
	{
		//writeKeyToJson("lab2/elgamal/resources/keys.json", Keys{Pair{A.a, A.k}, Pair{B.x, B.y}})
		encrypt(&A, &B, P)
		decrypt(&A, &B, P)
		//file, fileSize := methods.ReadFile("lab2/resourcesGlobal/test.jpg")
		//EncryptMessage(file, fileSize, &A, &B, P)
		//DecryptMessage(&A, &B, P)
		//methods.WriteFile("lab2/elgamal/resources/decrypt.jpg", B.decodedMessage)
	}

	fmt.Printf("Message of Alice: %d\nDecoded message of Bob: %d\n", A.message, B.decodedMessage)
}
