package blocks

import "fmt"
import "os"

type FileOutput struct {
	BlockData
	file *os.File
}

func (b *FileOutput) Update() {
	b.out = b.in
	b.file.Seek(0, 0)
	for _, v := range b.out {
		fmt.Fprintln(b.file, v)
	}
}

func ConstructFileOutput(x []string) Block {
	var file *os.File
	if x[0] == "stdout" {
		file = os.Stdout
	} else if x[0] == "stderr" {
		file = os.Stderr
	} else {
		file, _ = os.Create(x[0])
	}

	b := &FileOutput{file: file}
	return b
}

var ConstructFileOutputOk = AddConstructor("FileOutput", ConstructFileOutput)
