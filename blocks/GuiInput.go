package blocks

import (
	"../gui/"
	"../logger/"
	"../parser/"
	"fmt"
	"os"
	"strconv"
)

type GuiInput struct {
	InputBlockData
	v     []float64
	name  string
	fname string
	file  *os.File
}

func (b *GuiInput) Update() {
	b.out = b.v
	b.in = b.out
}

func (b *GuiInput) saveState() {
	b.file.Seek(0, 0)

	for _, value := range b.v {
		fmt.Fprintln(b.file, value)
	}
}

func (b *GuiInput) Put(x []float64) {
	if len(x) > 0 {
		b.v = x

		// also put into state file
		b.saveState()
	}

	b.in = x
	b.out = x
}

func (b *GuiInput) loadState() {
	file, errOpen := os.Open(b.fname)

	if errOpen == nil {
		vs := parser.VectorizeFileFloats(b.fname)

		b.v = vs
	}

	file.Close()
}

func GuiInputConstructor(name string, words []string) Block {
	values := []float64{}

	for _, word := range words {
		value, err := strconv.ParseFloat(word, 64)

		if err != nil {
			logger.WriteError("GuiInputConstructor()", err)
		}

		values = append(values, value)
	}

	// now try loading the stateFile if it exists
	fname := name + "_state.dat"

	b := &GuiInput{
		v:     values,
		name:  name,
		fname: fname,
	}

	b.loadState()

	var errCreate error
	b.file, errCreate = os.Create(fname)
	logger.WriteError("GuiInputConstructor()", errCreate)

	// register the Gui block
	gui.GuiBlocks[name] = b

	return b
}

var GuiInputConstructorOk = AddConstructor("GuiInput", GuiInputConstructor)
