package lines

import (
	"../blocks/"
	"../logger/"
	//"errors"
	"strconv"
)

type SplitLine struct {
	LineData
	nf int // number of floats per output
}

func (l *SplitLine) transfer() {
	x := l.b0[0].Get()

	for i, v := range l.b1 {
		i0 := i * l.nf
		i1 := (i + 1) * l.nf

		//fmt.Println(i0, i1, len(b.in), len(b.b1))
		if i0 > len(x)-1 {
			logger.WriteEvent("warning, SplitLine.Update(): too few output blocks")
			break
		}

		if i1 > len(x) {
			i1 = len(x)
			logger.WriteEvent("warning, SplitLine.Update(): too few output blocks, truncating")
		}

		v.Put(x[i0:i1])
	}
}

func SplitLineConstructor(name string, words []string, b map[string]blocks.Block) Line {
	nf, errInt := strconv.ParseInt(words[0], 10, 64)
	logger.WriteError("SplitLineConstructor()", errInt)

	b0 := getBlock(b, words[1])
	b1 := getBlocks(b, words[2:])

	l := &SplitLine{
		LineData: LineData{
			b0:        []blocks.Block{b0},
			b1:        b1,
			DebugName: name,
		},
		nf: int(nf),
	}

	return l
}

var SplitLineConstructorOk = AddConstructor("SplitLine", SplitLineConstructor)
