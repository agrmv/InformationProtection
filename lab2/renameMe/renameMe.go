package renameMe

import "os"

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func ReadFile(path string) ([]byte, int) {
	file, err := os.Open(path)
	checkError(err)

	defer func() {
		err := file.Close()
		checkError(err)
	}()

	fi, err := file.Stat()
	checkError(err)

	bytes := make([]byte, fi.Size())
	countBytes, err := file.Read(bytes)
	checkError(err)

	return bytes, countBytes
}

func WriteFile(bytes []byte) {
	file, err := os.Create("lab2/vernamCipher/resources/decode.txt")
	checkError(err)

	_, err = file.Write(bytes)
	checkError(err)
}
