package blocks

import (
  "../logger/"
	"bufio"
  //"fmt"
	"os"
	"strconv"
)

type FileInput struct {
	InputBlockData
	file    *os.File // usefull for seeking
	//scanner *bufio.Scanner
}

// column format is useful for reading
// row format is useful for piping
func (b *FileInput) Update() {
	b.file.Seek(0, 0) // also seeks to 0 when piping stdin

	b.out = []float64{}

  // rescan the file
	scanner := bufio.NewScanner(b.file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		x, _ := strconv.ParseFloat(scanner.Text(), 64)
    //fmt.Println("FileInput found: ", x)
		b.out = append(b.out, x)
	}

	b.in = b.out
}

func FileInputConstructor(name string, words []string) Block {
	var file *os.File
	if words[0] == "stdin" {
		file = os.Stdin
	} else {
    var errOpen error
		file, errOpen = os.Open(words[0])
    logger.WriteError("FileInputConstructor()", errOpen)
	}

	b := &FileInput{file: file}
	return b
}

var FileInputConstructorOk = AddConstructor("FileInput", FileInputConstructor)
