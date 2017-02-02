package blocks

import (
	"../logger/"
	"../parser/"
	"errors"
)

type ReductionNode struct {
	BlockData
	fn func([]float64) float64
}

func (b *ReductionNode) Put(x []float64) {
	b.in = x

	// output always has length 1
	if len(b.out) != 1 {
		b.out = make([]float64, 1)
	}

	b.out[0] = b.fn(x)
}

func ReductionNodeConstructor(name string, words []string) Block {
	var operator string
	positional := parser.PositionalArgs(&operator)
	parser.ParsePositionalArgs(words, positional)

	fns := make(map[string](func([]float64) float64))
	fns["And"] = ReductionNodeAndOperator
	fns["Or"] = ReductionNodeOrOperator

	var fn func([]float64) float64
	var ok bool
	if fn, ok = fns[operator]; !ok {
		validOperators := ""
		for key, _ := range fns {
			validOperators = validOperators + ", " + key
		}

		logger.WriteFatal("ReductionNodeConstructor()", errors.New("operator "+operator+" not found in \""+validOperators+"\""))
		return nil
	}

	b := &ReductionNode{
		fn: fn,
	}
	return b
}

// only 0.0 and 1.0 are treated, all other numbers are ignored
func ReductionNodeAndOperator(x []float64) float64 {
	y := 1.0

	for _, v := range x {
		if v == 0.0 {
			y = 0.0
			break
		}
	}

	return y
}

func ReductionNodeOrOperator(x []float64) float64 {
	y := 0.0

	for _, v := range x {
		if v == 1.0 {
			y = 1.0
			break
		}
	}

	return y
}
