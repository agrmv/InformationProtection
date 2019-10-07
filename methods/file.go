package methods

import (
	"os"
)

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func ReadFile(path string) ([]byte, int64) {
	file, err := os.Open(path)
	checkError(err)

	defer func() {
		err := file.Close()
		checkError(err)
	}()

	fileInfo, err := file.Stat()
	checkError(err)

	bytes := make([]byte, fileInfo.Size())
	_, err = file.Read(bytes)
	checkError(err)

	return bytes, fileInfo.Size()
}

func WriteFile(path string, bytes []byte) {
	file, err := os.Create(path)
	checkError(err)

	_, err = file.Write(bytes)
	checkError(err)
}
