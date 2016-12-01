package lines

import (
	"../blocks/"
)

type JoinLine struct {
	LineData
}

func (l *JoinLine) Transfer() {
	x := []float64{}

	for _, v := range l.b0 {
		x = append(x, blocks.Blocks[v].Get()...)
	}
	blocks.Blocks[l.b1[0]].Put(x)
}

func JoinLineConstructor(words []string) Line {
	b0 := blocks.checkNames(words[1:])
	b1 := blocks.checkName(words[0])

	l := &JoinLine{b0: b0, b1: []string{b1}}
	return l
}

var JoinLineConstructorOk = AddConstructor("JoinLine", JoinLineConstructor)
