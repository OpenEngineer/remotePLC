package blocks

import "strconv"

type ConstantInput struct {
	BlockData
	constant float64
}

func (b *ConstantInput) Update() {
	b.out = []float64{b.constant}
}

func ConstructConstantInput(words []string) Block {
	constant, _ := strconv.ParseFloat(words[0], 64)

	b := &ConstantInput{constant: constant}
	return b
}

var ConstantInputOk = AddConstructor("ConstantInput", ConstructConstantInput)
