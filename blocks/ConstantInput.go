package blocks

import (
	"log"
	"strconv"
)

type ConstantInput struct {
	InputBlockData
	constants []float64
}

func (b *ConstantInput) Update() {
	b.out = b.constants
	b.in = b.out
}

func ConstantInputConstructor(name string, words []string) Block {
	constants := []float64{}
	for _, word := range words {
		constant, err := strconv.ParseFloat(word, 64)
		if err != nil {
			log.Fatal(err)
		}

		constants = append(constants, constant)
	}

	b := &ConstantInput{constants: constants}
	return b
}

var ConstantInputConstructorOk = AddConstructor("ConstantInput", ConstantInputConstructor)
