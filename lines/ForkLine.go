package lines

import (
	"../blocks/"
)

type ForkLine struct {
	LineData
}

func (l *ForkLine) Transfer() {
	x := blocks.Blocks[l.b0[0]].Get()

	for _, v := range l.b1 {
		blocks.Blocks[v].Put(x)
	}
}

func ForkLineConstructor(words []string) Line {
	b0 := blocks.checkName(words[0])
	b1 := blocks.checkNames(words[1:])

	l := &ForkLine{b0: []string{b0}, b1: b1}
	return l
}

var ForkLineConstructorOk = AddConstructor("ForkLine", ForkLineConstructor)
