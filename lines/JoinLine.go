package lines

import (
	"../blocks/"
)

type JoinLine struct {
	LineData
}

func (l *JoinLine) transfer() {
	x := []float64{}

	for _, v := range l.b0 {
		x = append(x, v.Get()...)
	}
	l.b1[0].Put(x)
}

func JoinLineConstructor(words []string, b map[string]blocks.Block) Line {
	b0 := getBlocks(b, words[1:])
	b1 := getBlock(b, words[0])

	l := &JoinLine{
		LineData{
			b0:        b0,
			b1:        []blocks.Block{b1},
			DebugName: getDebugName("JoinLine", words),
		},
	}

	return l
}

var JoinLineConstructorOk = AddConstructor("JoinLine", JoinLineConstructor)
