package blocks

import (
)

type DiffLine struct {
	BlockData
	b0 []string // only two inputs, second minus first (i.e. smaller larger out)
	b1 string
}

func (b *DiffLine) Update() {
	b.in = []float64{}

	x := Blocks[b.b0[0]].Get()
	y := Blocks[b.b0[1]].Get()

	var n int
	if len(x) < len(y) {
		n = len(x)
	} else {
		n = len(y)
	}

	for i := 0; i < n; i++ {
		b.in = append(b.in, y[i]-x[i])
	}

	Blocks[b.b1].Put(b.in)

	b.out = b.in
}

func DiffLineConstructor(name string, words []string) Block {
	b0 := words[0:2]
	b1 := words[2]

	b := &DiffLine{b0: b0, b1: b1}
	return b
}

var DiffLineConstructorOk = AddConstructor("DiffLine", DiffLineConstructor)
