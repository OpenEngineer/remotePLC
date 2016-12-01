package blocks

import (
	"../logger/"
	//"errors"
	"fmt"
	"strconv"
)

type SplitLine struct {
	BlockData
	nf int // number of floats per output
	b0 string
	b1 []string
}

func (b *SplitLine) Update() {
	b.in = Blocks[b.b0].Get()
	//fmt.Println(b.in)

	for i, v := range b.b1 {
		if v == "_" {
			continue
		}

		if BlockMode == CONNECTIVITY {
			Blocks[v].Put([]float64{1})
		} else {
			i0 := i * b.nf
			i1 := (i + 1) * b.nf

			//fmt.Println(i0, i1, len(b.in), len(b.b1))
			if i0 > len(b.in)-1 {
				logger.WriteEvent("warning, SplitLine.Update(): too few output blocks")
				break
			}

			if i1 > len(b.in) {
				i1 = len(b.in)
				logger.WriteEvent("warning, SplitLine.Update(): too few output blocks, truncating")
			}
			Blocks[v].Put(b.in[i0:i1])
		}
	}

	b.out = b.in
}

func SplitLineConstructor(name string, words []string) Block {
	nf, errInt := strconv.ParseInt(words[0], 10, 64)
	logger.WriteError("SplitLineConstructor()", errInt)

	b0 := checkName(words[1])
	b1_ := words[2:]

	b1 := []string{}

	for _, v := range b1_ {
		if v == "_" {
      b1 = append(b1, "_") // "_" carries the info about it being a throwaway variable
		} else {
		  b1 = append(b1, checkName(v))
    }
	}
  fmt.Println(b0, b1)
	b := &SplitLine{nf: int(nf), b0: b0, b1: b1}
	return b
}

var SplitLineConstructorOk = AddConstructor("SplitLine", SplitLineConstructor)
