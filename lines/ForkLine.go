package lines

import (
	"../blocks/"
)

type ForkLine struct {
	LineData
}

func (l *ForkLine) Transfer() {
  if l.check() {
    x := l.b0[0].Get()

    for _, v := range l.b1 {
      v.Put(x)
    }
  }
}

func ForkLineConstructor(name string, words []string, b map[string]blocks.Block) Line {
	b0 := getBlock(b, words[0])
	b1 := getBlocks(b, words[1:])

	l := &ForkLine{
		LineData{
			b0:        []blocks.Block{b0},
			b1:        b1,
			DebugName: name,
		},
	}
	return l
}

var ForkLineConstructorOk = AddConstructor("ForkLine", ForkLineConstructor)
