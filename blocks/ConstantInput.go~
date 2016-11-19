package blocks

import "strconv"

type ConstantInput struct {
	InputBlockData
	constant float64
}

func (b *ConstantInput) Update() {
	b.out = []float64{b.constant}
	b.in = b.out
}

func ConstantInputConstructor(words []string) Block {
	constant, _ := strconv.ParseFloat(words[0], 64)

	b := &ConstantInput{constant: constant}
	return b
}

var ConstantInputConstructorOk = AddConstructor("ConstantInput", ConstantInputConstructor)
