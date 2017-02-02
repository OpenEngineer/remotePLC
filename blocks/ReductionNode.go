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

	fn := GetReductionNodeOperator(operator)

	b := &ReductionNode{
		fn: fn,
	}
	return b
}

var ReductionNodeConstructorOk = AddConstructor("ReductionNode", ReductionNodeConstructor)

var ReductionNodeOperators map[string](func([]float64) float64)

func AddReductionNodeOperator(operator string, fn func([]float64) float64) bool {
	var ok bool

	if _, ok := ReductionNodeOperators[operator]; !ok {
		logger.WriteError("AddReductionNodeOperator()", errors.New("operator already exists"))
	} else {
		ReductionNodeOperators[operator] = fn
	}

	return ok
}

func GetReductionNodeOperator(operator string) func([]float64) float64 {
	if fn, ok := ReductionNodeOperators[operator]; ok {
		return fn
	} else {
		validOperators := ""
		for key, _ := range ReductionNodeOperators {
			validOperators = validOperators + ", " + key
		}

		logger.WriteFatal("GetReductionNodeOperator()", errors.New("operator "+operator+" not found in \""+validOperators+"\""))
		return nil
	}
}
