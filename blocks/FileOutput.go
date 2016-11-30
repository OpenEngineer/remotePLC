package blocks

import "fmt"
import "os"

type FileOutput struct {
	OutputBlockData
	file *os.File
}

func (b *FileOutput) Update() {
	b.file.Seek(0, 0) // also seeks to 0 when piping stdout or stderr
	for _, v := range b.in {
		fmt.Fprintln(b.file, v)
	}
}

func FileOutputConstructor(name string, words []string) Block {
	var file *os.File
	if words[0] == "stdout" {
		file = os.Stdout
	} else if words[0] == "stderr" {
		file = os.Stderr
	} else {
		file, _ = os.Create(words[0])
	}

	b := &FileOutput{file: file}
	return b
}

var FileOutputConstructorOk = AddConstructor("FileOutput", FileOutputConstructor)
