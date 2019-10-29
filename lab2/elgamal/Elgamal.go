package main

import (
	"../../methods"
	"encoding/binary"
	"github.com/DenisDyachkov/i_p_labs/basic"
	dh "github.com/DenisDyachkov/i_p_labs/basic/diffie-hellman"
	"github.com/DenisDyachkov/i_p_labs/crypt/el_gamal"
	"io/ioutil"
	"math/rand"
	"os"
)

func encrypt(m byte, sessionKey int64, key *dh.Key) int64 {
	params := key.Params
	return (int64(m) * basic.FastPowByModule(key.PublicKey, sessionKey, params.P)) % params.P
}

func decrypt(r, e int64, key *dh.Key) int64 {
	params := key.Params
	power := params.P - 1 - key.PrivateKey
	return (e * basic.FastPowByModule(r, power, params.P)) % params.P
}

func SecretSessionKey(key *dh.Key) int64 {
	return rand.Int63n(key.Params.P-1) + 1
}

func EncryptMessage(message []byte, key *dh.Key) (addition []byte, encrypted []byte) {
	bin := make([]byte, 8)
	for _, m := range message {
		k := SecretSessionKey(key)
		a := basic.FastPowByModule(key.Params.G, k, key.Params.P)
		b := encrypt(m, k, key)
		binary.LittleEndian.PutUint64(bin, uint64(b))
		encrypted = append(encrypted, bin...)
		binary.LittleEndian.PutUint64(bin, uint64(a))
		addition = append(addition, bin...)
	}
	return
}

func DecryptMessage(r []byte, e []byte, key *dh.Key) []byte {
	var decrypted []byte
	for i := 0; i < len(e); i += 8 {
		_r := binary.LittleEndian.Uint64(r[i : i+8])
		_e := binary.LittleEndian.Uint64(e[i : i+8])
		c := decrypt(int64(_r), int64(_e), key)
		decrypted = append(decrypted, byte(c))
	}
	return decrypted
}

func main() {
	key := el_gamal.GenerateKey()
	src, _ := methods.ReadFile("lab2/resourcesGlobal/test.jpg")
	add, enc := el_gamal.EncryptMessage(src, key)
	dec := el_gamal.DecryptMessage(add, enc, key)
	el_gamal.SaveKeyToFile("key.el_gamal", key)
	//_ = ioutil.WriteFile("lab2/resourcesGlobal/test.jpg", src, os.ModePerm)
	_ = ioutil.WriteFile("lab2/elgamal/resources/elgamal-encrypted.jpg", enc, os.ModePerm)
	_ = ioutil.WriteFile("lab2/elgamal/resources/elgamal-encrypted-add.jpg", add, os.ModePerm)
	_ = ioutil.WriteFile("lab2/elgamal/resources/elgamal-decrypted.jpg", dec, os.ModePerm)

}
