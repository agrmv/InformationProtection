package main

import "fmt"

//хз как назвать
//проверяет на то принадлежит ли символ алфивиту удаляет пробелы и прочее и делает из всех юукв заглавные
func CheckClear(in *string) {
	var out []rune
	for _, v := range *in {
		if 65 <= v && v <= 90 {
			out = append(out, v)
		} else if 97 <= v && v <= 122 {
			out = append(out, v-32)
		}
	}
	*in = string(out)
}

func EncodePair(a, b rune) rune {
	return (((a - 'A') + (b - 'A')) % 26) + 'A'
}

func DecodePair(a, b rune) rune {
	return ((((a - 'A') - (b - 'A')) + 26) % 26) + 'A'
}

func Encode(msg, key string) string {
	CheckClear(&msg)
	CheckClear(&key)
	out := make([]rune, 0, len(msg))
	for i, v := range msg {
		out = append(out, EncodePair(v, rune(key[i%len(key)])))
	}
	return string(out)
}

func Decode(msg, key string) string {
	CheckClear(&msg)
	CheckClear(&key)
	out := make([]rune, 0, len(msg))
	for i, v := range msg {
		out = append(out, DecodePair(v, rune(key[i%len(key)])))
	}
	return string(out)
}

func main() {
	key := "papich"
	message := "Vsem moim bratyam salam"

	encoded := Encode(message, key)
	decoded := Decode(encoded, key)

	fmt.Println(encoded)
	fmt.Println(decoded)
}
