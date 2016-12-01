package lines

import (
	"../blocks/"
)

// only two inputs, second minus first (i.e. smaller larger out)
type DiffLine struct {
	LineData
}

func (l *DiffLine) Transfer() {
	x := blocks.Blocks[l.b0[0]].Get()
	y := blocks.Blocks[l.b0[1]].Get()

	var n int
	if len(x) < len(y) {
		n = len(x)
	} else {
		n = len(y)
	}

	d := []float64{}
	for i := 0; i < n; i++ {
		d = append(d, y[i]-x[i])
	}

	blocks.Blocks[l.b1[0]].Put(d)
}

func DiffLineConstructor(words []string) Line {
	b0 := blocks.checkNames(words[0:2])
	b1 := blocks.checkName(words[2])

	l := &DiffLine{b0: b0, b1: []string{b1}}
	return l
}

var DiffLineConstructorOk = AddConstructor("DiffLine", DiffLineConstructor)
