package lines

import (
	"../blocks/"
)

type JoinLine struct {
	LineData
}

func (l *JoinLine) Transfer() {
  if l.check() {
    x := []float64{}

    for _, v := range l.b0 {
      x = append(x, v.Get()...)
    }
    l.b1[0].Put(x)
  }
}

func JoinLineConstructor(name string, words []string, b map[string]blocks.Block) Line {
	b0 := getBlocks(b, words[1:])
	b1 := getBlock(b, words[0])

	l := &JoinLine{
		LineData{
			b0:        b0,
			b1:        []blocks.Block{b1},
			DebugName: name,
		},
	}

	return l
}

var JoinLineConstructorOk = AddConstructor("JoinLine", JoinLineConstructor)
