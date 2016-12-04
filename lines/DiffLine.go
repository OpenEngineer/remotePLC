package lines

import (
	"../blocks/"
)

// only two inputs, second minus first (i.e. smaller larger out)
type DiffLine struct {
	LineData
}

func (l *DiffLine) transfer() {
	x := l.b0[0].Get()
	y := l.b0[1].Get()

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

	l.b1[0].Put(d)
}

func DiffLineConstructor(name string, words []string, b map[string]blocks.Block) Line {
	b0 := getBlocks(b, words[0:2])
	b1 := getBlock(b, words[2])

	l := &DiffLine{
		LineData{
			b0:        b0,
			b1:        []blocks.Block{b1},
			DebugName: name,
		},
	}

	return l
}

var DiffLineConstructorOk = AddConstructor("DiffLine", DiffLineConstructor)
