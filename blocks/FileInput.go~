package blocks

import (
	"bufio"
	"os"
	"strconv"
)

type FileInput struct {
	InputBlockData
	file    *os.File // usefull for seeking
	scanner *bufio.Scanner
}

func (b *FileInput) Update() {
	b.file.Seek(0, 0) // also seeks to 0 when piping stdin

	b.out = []float64{}

	for b.scanner.Scan() {
		x, _ := strconv.ParseFloat(b.scanner.Text(), 64)
		b.out = append(b.out, x)
	}

	b.in = b.out
}

func FileInputConstructor(words []string) Block {
	var file *os.File
	if words[0] == "stdin" {
		file = os.Stdin
	} else {
		file, _ = os.Create(words[0])
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	b := &FileInput{file: file, scanner: scanner}
	return b
}

var FileInputConstructorOk = AddConstructor("FileInput", FileInputConstructor)
