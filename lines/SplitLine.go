package blocks

import (
	"../blocks/"
	"../logger/"
	//"errors"
	"fmt"
	"strconv"
)

type SplitLine struct {
	LineData
	nf int // number of floats per output
}

func (l *SplitLine) Transfer() {
	x := blocks.Blocks[l.b0[0]].Get()
	//fmt.Println(b.in)

	for i, v := range l.b1 {
		if v == "_" {
			continue
		}

		if blocks.BlockMode == CONNECTIVITY {
			blocks.Blocks[v].Put([]float64{1})
		} else {
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
			blocks.Blocks[v].Put(x[i0:i1])
		}
	}
}

func SplitLineConstructor(words []string) Line {
	nf, errInt := strconv.ParseInt(words[0], 10, 64)
	logger.WriteError("SplitLineConstructor()", errInt)

	b0 := blocks.checkName(words[1])
	b1_ := words[2:]

	b1 := []string{}

	for _, v := range b1_ {
		if v == "_" {
			b1 = append(b1, "_") // "_" carries the info about it being a throwaway variable
		} else {
			b1 = append(b1, blocks.checkName(v))
		}
	}
	//fmt.Println(b0, b1)
	l := &SplitLine{nf: int(nf), b0: []string{b0}, b1: b1}
	return l
}

var SplitLineConstructorOk = AddConstructor("SplitLine", SplitLineConstructor)
